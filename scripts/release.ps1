# Full Release Script for VPN Client

param(
    [Parameter(Mandatory=$true)][string]$Version,
    [switch]$NoBuild,
    [switch]$NoPackage,
    [switch]$DryRun,
    [switch]$Help
)

Write-Host "Release script for VPN Client"
Write-Host "Version: $Version"

if ($Help) {
    Write-Host "Usage: release.ps1 <version> [options]"
    Write-Host "Options:"
    Write-Host "  -NoBuild     Skip building binaries"
    Write-Host "  -NoPackage   Skip packaging"
    Write-Host "  -DryRun      Show what would be done without doing it"
    Write-Host "  -Help        Display this help message"
    exit 0
}

# Validate version format (should be vX.Y.Z)
if ($Version -notmatch "^v\d+\.\d+\.\d+$") {
    Write-Host "Error: Version must be in format vX.Y.Z (e.g., v1.0.0)" -ForegroundColor Red
    exit 1
}

Write-Host "Preparing to release VPN Client $Version" -ForegroundColor Blue

# Check if working directory is clean
$gitStatus = git status --porcelain
if ($gitStatus) {
    Write-Host "Warning: Working directory is not clean" -ForegroundColor Yellow
    if (-not $DryRun) {
        $response = Read-Host "Continue anyway? (y/N)"
        if ($response -notmatch "^[Yy]$") {
            exit 1
        }
    }
}

# Update version in versioninfo.json
if (Test-Path "versioninfo.json") {
    Write-Host "Updating version in versioninfo.json..." -ForegroundColor Yellow
    if (-not $DryRun) {
        # Extract version numbers
        $versionParts = $Version.Substring(1).Split(".")
        $major = $versionParts[0]
        $minor = $versionParts[1]
        $patch = $versionParts[2]
        
        # Read the JSON file
        $versionInfo = Get-Content "versioninfo.json" | ConvertFrom-Json
        
        # Update version fields
        $versionInfo.FixedFileInfo.FileVersion.Major = [int]$major
        $versionInfo.FixedFileInfo.FileVersion.Minor = [int]$minor
        $versionInfo.FixedFileInfo.FileVersion.Patch = [int]$patch
        $versionInfo.FixedFileInfo.ProductVersion.Major = [int]$major
        $versionInfo.FixedFileInfo.ProductVersion.Minor = [int]$minor
        $versionInfo.FixedFileInfo.ProductVersion.Patch = [int]$patch
        $versionInfo.StringFileInfo.FileVersion = "$major.$minor.$patch.0"
        $versionInfo.StringFileInfo.ProductVersion = "$major.$minor.$patch.0"
        
        # Save the updated JSON
        $versionInfo | ConvertTo-Json -Depth 3 | Set-Content "versioninfo.json"
    }
    Write-Host "OK Updated versioninfo.json" -ForegroundColor Green
}

# Update version in NSIS script (installer.nsi) or Inno Setup script (installer.iss)
$installerUpdated = $false
if (Test-Path "installer.nsi") {
    Write-Host "Updating version in installer.nsi..." -ForegroundColor Yellow
    if (-not $DryRun) {
        $content = Get-Content "installer.nsi"
        $content = $content -replace '!define APP_VERSION ".*"', "!define APP_VERSION `"$($Version.Substring(1))`""
        Set-Content "installer.nsi" $content
    }
    Write-Host "OK Updated installer.nsi" -ForegroundColor Green
    $installerUpdated = $true
}

# Then check for Inno Setup script
if (Test-Path "installer.iss") {
    Write-Host "Updating version in installer.iss..." -ForegroundColor Yellow
    if (-not $DryRun) {
        $content = Get-Content "installer.iss"
        $content = $content -replace '#define MyAppVersion ".*"', "#define MyAppVersion `"$($Version.Substring(1))`""
        Set-Content "installer.iss" $content
    }
    Write-Host "OK Updated installer.iss" -ForegroundColor Green
    $installerUpdated = $true
}

if (-not $installerUpdated) {
    Write-Host "Warning: No installer script found (neither installer.nsi nor installer.iss)" -ForegroundColor Yellow
}

# Build binaries if requested
if (-not $NoBuild) {
    Write-Host "Building binaries..." -ForegroundColor Yellow
    if (-not $DryRun) {
        powershell -ExecutionPolicy Bypass -File ".\scripts\build-all.ps1" $($Version.Substring(1))
    } else {
        Write-Host "[DRY RUN] Would run: .\scripts\build-all.ps1 $($Version.Substring(1))" -ForegroundColor Blue
    }
} else {
    Write-Host "Skipping build phase" -ForegroundColor Yellow
}

# Build Android app if requested
if (-not $NoBuild) {
    Write-Host "Building Android packages..." -ForegroundColor Yellow
    if (-not $DryRun) {
        # Check if Flutter project exists
        if (Test-Path "ui\mobile\flutter_vpn") {
            Write-Host "Building Flutter Android app..." -ForegroundColor Yellow
            Set-Location "ui\mobile\flutter_vpn"
            
            # Build APKs
            flutter build apk --release --split-per-abi
            # Build App Bundle
            flutter build appbundle --release
            
            Set-Location "..\..\.."
        } else {
            Write-Host "Warning: Flutter project not found, skipping Android build" -ForegroundColor Yellow
        }
    } else {
        Write-Host "[DRY RUN] Would build Android app with Flutter" -ForegroundColor Blue
    }
} else {
    Write-Host "Skipping Android build phase" -ForegroundColor Yellow
}

# Package if requested
if (-not $NoPackage) {
    Write-Host "Creating packages..." -ForegroundColor Yellow
    if (-not $DryRun) {
        powershell -ExecutionPolicy Bypass -File ".\scripts\package.ps1" $($Version.Substring(1))
        
        # Package Android apps
        if (Test-Path "ui\mobile\flutter_vpn") {
            Write-Host "Packaging Android apps..." -ForegroundColor Yellow
            # Copy Android build outputs to packages directory
            $androidBuildPath = "ui\mobile\flutter_vpn\build\app\outputs"
            if (Test-Path $androidBuildPath) {
                # Create android directory in packages
                New-Item -ItemType Directory -Force -Path "packages\android" | Out-Null
                
                # Copy APKs
                if (Test-Path "$androidBuildPath\flutter-apk") {
                    Get-ChildItem "$androidBuildPath\flutter-apk\app-*-release.apk" | ForEach-Object {
                        Copy-Item $_.FullName "packages\android\vpn-client-$($Version.Substring(1))-$(($_.Name -replace 'app-', '') -replace '-release.apk', '.apk')"
                    }
                }
                
                # Copy App Bundle
                if (Test-Path "$androidBuildPath\bundle\release\app-release.aab") {
                    Copy-Item "$androidBuildPath\bundle\release\app-release.aab" "packages\android\vpn-client-$($Version.Substring(1))-appbundle.aab"
                }
                
                Write-Host "OK Android packages created" -ForegroundColor Green
            }
        }
    } else {
        Write-Host "[DRY RUN] Would run: .\scripts\package.ps1 $($Version.Substring(1))" -ForegroundColor Blue
        Write-Host "[DRY RUN] Would package Android apps" -ForegroundColor Blue
    }
} else {
    Write-Host "Skipping packaging phase" -ForegroundColor Yellow
}

# Create Git tag
Write-Host "Creating Git tag..." -ForegroundColor Yellow
if (-not $DryRun) {
    git tag -a "$Version" -m "Release $Version"
} else {
    Write-Host "[DRY RUN] Would create Git tag: $Version" -ForegroundColor Blue
}
Write-Host "OK Tag created" -ForegroundColor Green

# Summary
Write-Host ""
Write-Host "==================================================" -ForegroundColor Green
Write-Host "           RELEASE PREPARATION COMPLETE" -ForegroundColor Green
Write-Host "==================================================" -ForegroundColor Green
Write-Host "Version:          $Version" -ForegroundColor Green
Write-Host "Build performed:  $(-not $NoBuild)" -ForegroundColor Green
Write-Host "Packaging done:   $(-not $NoPackage)" -ForegroundColor Green
Write-Host "Git tag created:  Yes" -ForegroundColor Green
Write-Host "" -ForegroundColor Green
Write-Host "Next steps:" -ForegroundColor Green
Write-Host "1. Review the build artifacts in dist/ and packages/" -ForegroundColor Green
Write-Host "2. Test the packages on target platforms" -ForegroundColor Green
Write-Host "3. Push the tag: git push origin $Version" -ForegroundColor Green
Write-Host "4. Create GitHub release: gh release create $Version packages/*" -ForegroundColor Green
Write-Host "==================================================" -ForegroundColor Green
