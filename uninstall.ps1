$app = 'com.dannyvankooten.browserpass'
$dirpath = Join-Path -Path $env:localappdata -ChildPath 'browserpass'

if (Test-Path -Path $dirPath) {
    Remove-Item -Recurse -Force $dirpath
}

If (-NOT ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole(`
            [Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Warning "Please re-run this script with Admin rights!"
    exit
}

If (Test-Path -Path "hklm:\Software\Google\Chrome\NativeMessagingHosts\$app") {
    Write-Host "Uninstalling for Chrome - all users"
    Remove-Item -Path "hklm:\Software\Google\Chrome\NativeMessagingHosts\$app" -force
}
if (Test-Path -Path "hklm:\Software\Mozilla\NativeMessagingHosts\$app") {
    Write-Host "Uninstalling for Firefox - all users"
    Remove-Item -Path "hklm:\Software\Mozilla\NativeMessagingHosts\$app" -force
}

if (Test-Path -Path "hkcu:\Software\Mozilla\NativeMessagingHosts\$app") {
    Write-Host "Uninstalling for Firefox - local user"
    Remove-Item -Path "hkcu:\Software\Mozilla\NativeMessagingHosts\$app" -force
}

if (Test-Path -Path "hkcu:\Software\Google\Chrome\NativeMessagingHosts\$app") {
    Write-Host "Uninstalling for Chrome - local user"
    Remove-Item -Path "hkcu:\Software\Google\Chrome\NativeMessagingHosts\$app" -force
}