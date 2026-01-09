# Create Windows Scheduled Task: Run sync script every 4 hours
# Run this script as Administrator

$TaskName = "NewAPI-AutoSync"
$ScriptPath = "C:\Users\admin\Desktop\start\code\new-api\sync_newapi.ps1"

# Remove existing task
Unregister-ScheduledTask -TaskName $TaskName -Confirm:$false -ErrorAction SilentlyContinue

# Create trigger: every 4 hours
$Trigger = New-ScheduledTaskTrigger -Once -At (Get-Date) -RepetitionInterval (New-TimeSpan -Hours 4)

# Create action: run PowerShell script
$Action = New-ScheduledTaskAction -Execute "powershell.exe" -Argument "-ExecutionPolicy Bypass -File `"$ScriptPath`""

# Create settings
$Settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -StartWhenAvailable

# Register task
Register-ScheduledTask -TaskName $TaskName -Trigger $Trigger -Action $Action -Settings $Settings -Description "Auto sync and deploy new-api"

Write-Host "Task '$TaskName' created, runs every 4 hours" -ForegroundColor Green
Write-Host "View in Task Scheduler" -ForegroundColor Yellow
