#Requires -Version 5.1
<#
.SYNOPSIS
    FCF - Find File or Folder Command for Windows
    Interactive file/folder finder with parallel search and real-time streaming

.DESCRIPTION
    FCF is an interactive command-line tool for finding files and folders with
    advanced pattern matching, real-time streaming results, and easy navigation.
    It uses fd for parallel search when available, falling back to Get-ChildItem.

.PARAMETER Pattern
    The file/folder pattern to search for (supports wildcards)

.PARAMETER Path
    The directory to search in (defaults to current directory)

.PARAMETER IgnoreCase
    Perform case-insensitive search

.PARAMETER Type
    Filter by type: 'f' for files, 'd' for directories

.PARAMETER ShowSize
    Display file sizes in the output

.PARAMETER MaxDisplay
    Maximum number of results to display (0 = unlimited)

.PARAMETER Help
    Show help message

.EXAMPLE
    .\fcf.ps1
    Starts interactive mode

.EXAMPLE
    .\fcf.ps1 "*.log"
    Find all .log files in current directory

.EXAMPLE
    .\fcf.ps1 "*.js" C:\Projects
    Find all .js files in C:\Projects

.EXAMPLE
    .\fcf.ps1 -IgnoreCase "*.PNG"
    Case-insensitive search for PNG files
#>

[CmdletBinding()]
param(
    [Parameter(Position = 0)]
    [string]$Pattern,

    [Parameter(Position = 1)]
    [string]$Path,

    [Alias("i")]
    [switch]$IgnoreCase,

    [Alias("t")]
    [ValidateSet("f", "d")]
    [string]$Type,

    [switch]$ShowSize,

    [int]$MaxDisplay = 0,

    [Alias("h")]
    [switch]$Help
)

$Version = "1.0.0"

# Enable ANSI colors for Windows Terminal / PowerShell 7+
if ($PSVersionTable.PSVersion.Major -ge 7 -or $env:WT_SESSION) {
    $PSStyle.OutputRendering = 'Ansi'
}

# Color codes (ANSI)
$script:Colors = @{
    Red     = "`e[0;31m"
    Green   = "`e[0;32m"
    Yellow  = "`e[1;33m"
    Blue    = "`e[0;34m"
    Cyan    = "`e[0;36m"
    Magenta = "`e[0;35m"
    Bold    = "`e[1m"
    Dim     = "`e[2m"
    NC      = "`e[0m"
}

# Fallback for older PowerShell without ANSI support
if ($PSVersionTable.PSVersion.Major -lt 7 -and -not $env:WT_SESSION) {
    $script:Colors = @{
        Red     = ""
        Green   = ""
        Yellow  = ""
        Blue    = ""
        Cyan    = ""
        Magenta = ""
        Bold    = ""
        Dim     = ""
        NC      = ""
    }
}

# Store results for navigation
$script:Results = @()

# Temp file for navigation path
$script:NavPathFile = Join-Path $env:TEMP "fcf_nav_path"

# Clean up any stale navigation path on start
if (Test-Path $script:NavPathFile) {
    Remove-Item $script:NavPathFile -Force -ErrorAction SilentlyContinue
}

function Show-Help {
    $c = $script:Colors
    Write-Host "$($c.Bold)fcf - Find File or Folder$($c.NC) v$Version"
    Write-Host ""
    Write-Host "$($c.Bold)USAGE:$($c.NC)"
    Write-Host "    fcf [OPTIONS] [PATTERN] [PATH]"
    Write-Host "    fcf                          # Interactive mode"
    Write-Host ""
    Write-Host "$($c.Bold)DESCRIPTION:$($c.NC)"
    Write-Host "    Interactive tool to find files and folders with pattern matching"
    Write-Host "    and real-time streaming results. Uses parallel search for speed."
    Write-Host ""
    Write-Host "$($c.Bold)OPTIONS:$($c.NC)"
    Write-Host "    $($c.Cyan)-Help, -h$($c.NC)               Show this help message"
    Write-Host "    $($c.Cyan)-IgnoreCase, -i$($c.NC)         Case-insensitive pattern matching"
    Write-Host "    $($c.Cyan)-Type TYPE, -t TYPE$($c.NC)     Filter by type: $($c.Yellow)f$($c.NC)(file) or $($c.Yellow)d$($c.NC)(directory)"
    Write-Host "    $($c.Cyan)-ShowSize$($c.NC)               Display file sizes"
    Write-Host "    $($c.Cyan)-MaxDisplay NUM$($c.NC)         Maximum results to display (default: unlimited)"
    Write-Host ""
    Write-Host "$($c.Bold)SHELL INTEGRATION (for navigation to work):$($c.NC)"
    Write-Host "    Add this to your $($c.Cyan)`$PROFILE$($c.NC):"
    Write-Host ""
    Write-Host "    $($c.Yellow)function fcf {$($c.NC)"
    Write-Host "    $($c.Yellow)    & `"C:\path\to\fcf.ps1`" @args$($c.NC)"
    Write-Host "    $($c.Yellow)    `$navPath = `"`$env:TEMP\fcf_nav_path`"$($c.NC)"
    Write-Host "    $($c.Yellow)    if (Test-Path `$navPath) {$($c.NC)"
    Write-Host "    $($c.Yellow)        Set-Location (Get-Content `$navPath)$($c.NC)"
    Write-Host "    $($c.Yellow)        Remove-Item `$navPath -Force$($c.NC)"
    Write-Host "    $($c.Yellow)    }$($c.NC)"
    Write-Host "    $($c.Yellow)}$($c.NC)"
    Write-Host ""
    Write-Host "$($c.Bold)EXAMPLES:$($c.NC)"
    Write-Host "    $($c.Green)# Interactive mode$($c.NC)"
    Write-Host "    fcf"
    Write-Host ""
    Write-Host "    $($c.Green)# Find all .log files$($c.NC)"
    Write-Host "    fcf `"*.log`""
    Write-Host ""
    Write-Host "    $($c.Green)# Find in specific directory$($c.NC)"
    Write-Host "    fcf `"*.js`" C:\Projects"
    Write-Host ""
    Write-Host "    $($c.Green)# Case-insensitive search for PNG files$($c.NC)"
    Write-Host "    fcf -i `"*.PNG`""
    Write-Host ""
    Write-Host "    $($c.Green)# Find only directories named 'src'$($c.NC)"
    Write-Host "    fcf -t d src"
    Write-Host ""
    Write-Host "    $($c.Green)# Find with file sizes$($c.NC)"
    Write-Host "    fcf -ShowSize `"*.mp4`""
    Write-Host ""
    Write-Host "$($c.Bold)INTERACTIVE WORKFLOW:$($c.NC)"
    Write-Host "    Step 1: Enter path to search"
    Write-Host "    Step 2: Enter pattern to find"
    Write-Host "    Step 3: Navigate to a result path"
    Write-Host ""
    Write-Host "$($c.Bold)NAVIGATION OPTIONS:$($c.NC)"
    Write-Host "    After navigation, choose:"
    Write-Host "    $($c.Cyan)f$($c.NC) - Find again (restart from Step 1)"
    Write-Host "    $($c.Cyan)r$($c.NC) - Repeat search (go to Step 2, same path)"
    Write-Host "    $($c.Cyan)n$($c.NC) - Exit"
    Write-Host ""
    Write-Host "$($c.Bold)PERFORMANCE:$($c.NC)"
    Write-Host "    - Uses 'fd' for fast parallel searching (if installed)"
    Write-Host "    - Falls back to Get-ChildItem if fd is not available"
    Write-Host "    - Install fd: $($c.Cyan)winget install sharkdp.fd$($c.NC)"
    Write-Host "                  $($c.Cyan)choco install fd$($c.NC)"
    Write-Host "                  $($c.Cyan)scoop install fd$($c.NC)"
    Write-Host ""
}

function Test-FdInstalled {
    $null -ne (Get-Command fd -ErrorAction SilentlyContinue)
}

function Format-FileSize {
    param([long]$Bytes)

    if ($Bytes -ge 1GB) {
        return "{0:N1}G" -f ($Bytes / 1GB)
    }
    elseif ($Bytes -ge 1MB) {
        return "{0:N1}M" -f ($Bytes / 1MB)
    }
    elseif ($Bytes -ge 1KB) {
        return "{0:N1}K" -f ($Bytes / 1KB)
    }
    else {
        return "${Bytes}B"
    }
}

function Get-FileInfo {
    param([string]$FilePath)

    $info = ""
    if ($ShowSize -and (Test-Path $FilePath -PathType Leaf)) {
        $size = (Get-Item $FilePath).Length
        $info = " $($script:Colors.Dim)($(Format-FileSize $size))$($script:Colors.NC)"
    }
    return $info
}

function Show-Result {
    param(
        [string]$FilePath,
        [int]$Count
    )

    $c = $script:Colors
    $info = Get-FileInfo -FilePath $FilePath
    $item = Get-Item $FilePath -ErrorAction SilentlyContinue

    if ($item.PSIsContainer) {
        # Directory
        Write-Host "$($c.Cyan)  [$Count]$($c.NC) $($c.Blue)📁 $FilePath\$($c.NC)$info"
    }
    elseif ($item.Extension -in @('.exe', '.bat', '.cmd', '.ps1', '.com')) {
        # Executable
        Write-Host "$($c.Cyan)  [$Count]$($c.NC) $($c.Green)⚡ $FilePath$($c.NC)$info"
    }
    elseif ($item.Attributes -band [System.IO.FileAttributes]::ReparsePoint) {
        # Symlink/Junction
        Write-Host "$($c.Cyan)  [$Count]$($c.NC) $($c.Magenta)🔗 $FilePath$($c.NC)$info"
    }
    else {
        # Regular file
        Write-Host "$($c.Cyan)  [$Count]$($c.NC) 📄 $FilePath$info"
    }
}

function Show-Header {
    Clear-Host
    $c = $script:Colors
    Write-Host "$($c.Bold)$($c.Cyan)╔════════════════════════════════════════╗$($c.NC)"
    Write-Host "$($c.Bold)$($c.Cyan)║$($c.NC)   $($c.Bold)fcf$($c.NC) - Find File or Folder          $($c.Bold)$($c.Cyan)║$($c.NC)"
    Write-Host "$($c.Bold)$($c.Cyan)╚════════════════════════════════════════╝$($c.NC)"
    Write-Host ""
}

function Invoke-FdSearch {
    param(
        [string]$SearchPattern,
        [string]$SearchPath
    )

    $fdArgs = @("--color", "never", "--hidden", "--no-ignore")

    # Type filter
    if ($Type -eq "f") {
        $fdArgs += @("-t", "f")
    }
    elseif ($Type -eq "d") {
        $fdArgs += @("-t", "d")
    }

    # Case sensitivity
    if ($IgnoreCase) {
        $fdArgs += "-i"
    }
    else {
        $fdArgs += "-s"
    }

    # Glob pattern
    $fdArgs += @("-g", $SearchPattern, $SearchPath)

    & fd @fdArgs 2>$null
}

function Invoke-GetChildItemSearch {
    param(
        [string]$SearchPattern,
        [string]$SearchPath
    )

    $params = @{
        Path    = $SearchPath
        Recurse = $true
        Force   = $true  # Include hidden files
        ErrorAction = 'SilentlyContinue'
    }

    # Type filter
    if ($Type -eq "f") {
        $params['File'] = $true
    }
    elseif ($Type -eq "d") {
        $params['Directory'] = $true
    }

    # Get items and filter by pattern
    Get-ChildItem @params | Where-Object {
        if ($IgnoreCase) {
            $_.Name -like $SearchPattern
        }
        else {
            # Case-sensitive comparison using -clike
            $_.Name -clike $SearchPattern
        }
    } | ForEach-Object { $_.FullName }
}

function Navigate-ToPath {
    param([string]$TargetPath)

    $c = $script:Colors

    # If it's a file, get the directory
    if (Test-Path $TargetPath -PathType Leaf) {
        $TargetPath = Split-Path $TargetPath -Parent
    }

    # Check if path exists
    if (-not (Test-Path $TargetPath -PathType Container)) {
        Write-Host "$($c.Red)ERROR:$($c.NC) Directory '$TargetPath' does not exist"
        return $false
    }

    # Get absolute path
    $TargetPath = (Resolve-Path $TargetPath).Path

    # Save path to temp file for shell integration
    $TargetPath | Out-File -FilePath $script:NavPathFile -Encoding UTF8 -NoNewline

    Write-Host "$($c.Green)✓ Will navigate to:$($c.NC) $($c.Cyan)$TargetPath$($c.NC)"
    Write-Host ""
    Write-Host "$($c.Dim)Contents:$($c.NC)"
    Get-ChildItem $TargetPath | Format-Table -AutoSize
    return $true
}

function Show-OptionsMenu {
    $c = $script:Colors
    Write-Host ""
    Write-Host "$($c.Bold)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$($c.NC)"
    Write-Host "$($c.Bold)Options:$($c.NC)"
    Write-Host "  $($c.Cyan)[f]$($c.NC) Find again (new search)"
    Write-Host "  $($c.Cyan)[r]$($c.NC) Repeat search (same path)"
    Write-Host "  $($c.Cyan)[n]$($c.NC) Exit"
    Write-Host "$($c.Bold)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$($c.NC)"
    Write-Host ""

    $choice = Read-Host "$($c.Cyan)Choose$($c.NC)"

    switch ($choice.ToLower()) {
        "f" { return 1 }  # Go to Step 1
        "r" { return 2 }  # Go to Step 2
        default { return 0 }  # Exit
    }
}

# Main execution
if ($Help) {
    Show-Help
    exit 0
}

$CurrentStep = 1
$SearchPath = $Path
$SearchPattern = $Pattern
$CliMode = -not [string]::IsNullOrEmpty($Pattern)

while ($true) {
    # Reset results
    $script:Results = @()

    # Show header
    Show-Header

    $c = $script:Colors

    # Step 1: Get search path
    if ($CurrentStep -eq 1) {
        if ([string]::IsNullOrEmpty($SearchPath) -or -not $CliMode) {
            Write-Host "$($c.Bold)Step 1:$($c.NC) Enter path to search"
            Write-Host "$($c.Dim)(Press Enter for current directory: $PWD)$($c.NC)"
            Write-Host ""
            $userPath = Read-Host "$($c.Cyan)Path$($c.NC)"

            if ([string]::IsNullOrEmpty($userPath)) {
                $SearchPath = "."
                Write-Host "$($c.Green)Using current directory$($c.NC)"
            }
            else {
                # Expand environment variables and ~ to home directory
                $SearchPath = $userPath -replace "^~", $env:USERPROFILE
                $SearchPath = [Environment]::ExpandEnvironmentVariables($SearchPath)

                if (-not (Test-Path $SearchPath -PathType Container)) {
                    Write-Host "$($c.Red)ERROR:$($c.NC) Directory '$SearchPath' does not exist"
                    Read-Host "Press Enter to try again..."
                    continue
                }
            }
            Write-Host ""
        }
        $CurrentStep = 2
    }

    # Step 2: Get pattern
    if ($CurrentStep -eq 2) {
        if ([string]::IsNullOrEmpty($SearchPattern) -or -not $CliMode) {
            Write-Host "$($c.Bold)Step 2:$($c.NC) Enter file/folder name or pattern to find"
            Write-Host "$($c.Dim)Examples: *.log, config, .env, src, *.js$($c.NC)"
            Write-Host ""
            $SearchPattern = Read-Host "$($c.Cyan)Pattern$($c.NC)"

            if ([string]::IsNullOrEmpty($SearchPattern)) {
                Write-Host "$($c.Red)ERROR:$($c.NC) Pattern cannot be empty"
                Read-Host "Press Enter to try again..."
                continue
            }
            Write-Host ""
        }
    }

    # Show search info
    Write-Host "$($c.Bold)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$($c.NC)"
    Write-Host "$($c.Blue)Searching in:$($c.NC) $($c.Cyan)$SearchPath$($c.NC)"
    Write-Host "$($c.Blue)Pattern:$($c.NC) $($c.Yellow)$SearchPattern$($c.NC)"

    # Show search method
    if (Test-FdInstalled) {
        Write-Host "$($c.Blue)Method:$($c.NC) $($c.Green)fd (parallel search)$($c.NC)"
    }
    else {
        Write-Host "$($c.Blue)Method:$($c.NC) $($c.Yellow)Get-ChildItem (sequential - install 'fd' for faster search)$($c.NC)"
    }
    Write-Host "$($c.Bold)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$($c.NC)"
    Write-Host ""

    # Execute search with real-time streaming
    Write-Host "$($c.Bold)Results:$($c.NC) $($c.Dim)(streaming in real-time...)$($c.NC)"
    Write-Host ""

    $count = 0
    $stopwatch = [System.Diagnostics.Stopwatch]::StartNew()

    # Resolve search path
    $resolvedPath = if ($SearchPath -eq ".") { $PWD.Path } else { (Resolve-Path $SearchPath).Path }

    if (Test-FdInstalled) {
        $searchResults = Invoke-FdSearch -SearchPattern $SearchPattern -SearchPath $resolvedPath
    }
    else {
        $searchResults = Invoke-GetChildItemSearch -SearchPattern $SearchPattern -SearchPath $resolvedPath
    }

    foreach ($file in $searchResults) {
        if ([string]::IsNullOrEmpty($file)) { continue }

        $count++
        $script:Results += $file

        # Check max display limit
        if ($MaxDisplay -gt 0 -and $count -gt $MaxDisplay) {
            continue  # Still count but don't display
        }

        # Display result immediately (streaming)
        Show-Result -FilePath $file -Count $count
    }

    $stopwatch.Stop()
    $elapsed = "{0:N2}" -f $stopwatch.Elapsed.TotalSeconds

    # Summary
    Write-Host ""
    Write-Host "$($c.Bold)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$($c.NC)"

    if ($count -eq 0) {
        Write-Host "$($c.Yellow)No matches found$($c.NC) for pattern: $($c.Cyan)$SearchPattern$($c.NC)"
        Write-Host ""
        Write-Host "$($c.Dim)Tips:$($c.NC)"
        Write-Host "  - Try a different pattern"
        Write-Host "  - Use $($c.Cyan)-IgnoreCase$($c.NC) for case-insensitive search"

        # Show options even when no results
        $result = Show-OptionsMenu

        if ($result -eq 0) {
            exit 0
        }
        elseif ($result -eq 1) {
            $CurrentStep = 1
            $SearchPath = ""
            $SearchPattern = ""
            $CliMode = $false
            continue
        }
        elseif ($result -eq 2) {
            $CurrentStep = 2
            $SearchPattern = ""
            $CliMode = $false
            continue
        }
    }
    else {
        Write-Host "$($c.Green)$($c.Bold)Found $count match(es)$($c.NC) in $($c.Cyan)${elapsed}s$($c.NC)"

        if ($MaxDisplay -gt 0 -and $count -gt $MaxDisplay) {
            Write-Host "$($c.Yellow)(Displayed first $MaxDisplay of $count)$($c.NC)"
        }
    }

    Write-Host "$($c.Bold)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$($c.NC)"

    # Step 3: Navigate to path
    Write-Host ""
    Write-Host "$($c.Bold)Step 3:$($c.NC) Enter path to navigate to"
    Write-Host "$($c.Dim)(Enter a number from results, full path, or press Enter to skip)$($c.NC)"
    Write-Host ""
    $navInput = Read-Host "$($c.Cyan)Navigate to$($c.NC)"

    if (-not [string]::IsNullOrEmpty($navInput)) {
        # Check if input is a number (result index)
        if ($navInput -match '^\d+$') {
            $index = [int]$navInput - 1
            if ($index -ge 0 -and $index -lt $script:Results.Count) {
                $targetPath = $script:Results[$index]
            }
            else {
                Write-Host "$($c.Red)ERROR:$($c.NC) Invalid result number"
                $targetPath = $null
            }
        }
        else {
            # Treat as path (expand ~ and env vars)
            $targetPath = $navInput -replace "^~", $env:USERPROFILE
            $targetPath = [Environment]::ExpandEnvironmentVariables($targetPath)
        }

        if ($targetPath) {
            Write-Host ""
            $null = Navigate-ToPath -TargetPath $targetPath
        }
    }
    else {
        Write-Host "$($c.Dim)Skipped navigation$($c.NC)"
    }

    # Show options menu after Step 3
    $result = Show-OptionsMenu

    if ($result -eq 0) {
        Write-Host "$($c.Green)Goodbye!$($c.NC)"
        exit 0
    }
    elseif ($result -eq 1) {
        # Go to Step 1
        $CurrentStep = 1
        $SearchPath = ""
        $SearchPattern = ""
        $CliMode = $false
    }
    elseif ($result -eq 2) {
        # Go to Step 2
        $CurrentStep = 2
        $SearchPattern = ""
        $CliMode = $false
    }
}

exit 0
