; -- Inno Setup Script for TAMALabs with Windows Service Registration --
; ===================================================================
; FEATURES:
; - Automatic service installation and startup
; - System tray auto-start with administrator privileges (via scheduled task)
; - Avoid duplicate tray instances
; - Clean uninstallation (removes service, scheduled task, tray process)
; ===================================================================

[Setup]
AppId={{F4A4A2A2-702D-4B1F-A88E-5E3A1A8E2E8A}}
AppName=TAMALabs
AppVersion=1.0
AppPublisher=Elga Tama
AppPublisherURL=https://www.elgatama.com/
AppSupportURL=https://www.elgatama.com/support
AppUpdatesURL=https://www.elgatama.com/updates
DefaultDirName={autopf}\TAMALabs
DefaultGroupName=TAMALabs
AllowNoIcons=yes
OutputBaseFilename=TAMALabs-setup-v1.0
OutputDir=.\installers
Compression=lzma
SolidCompression=yes
WizardStyle=modern
ArchitecturesInstallIn64BitMode=x64compatible
PrivilegesRequired=admin
PrivilegesRequiredOverridesAllowed=dialog

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Dirs]
Name: "{app}\logs"; Permissions: users-modify


[Files]
Source: "bin\TAMALabs.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\nssm.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\TAMALabsTray.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\service-helper.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: ".env"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{autodesktop}\TAMALabs (Open Web)"; \
	Filename: "rundll32.exe"; \
	Parameters: "url.dll,FileProtocolHandler http://127.0.0.1:8322"; \
	IconFilename: "{app}\TAMALabs.exe"; \
	Tasks: desktopicon

; Fallback shortcut if scheduled task fails
Name: "{userstartup}\TAMALabs Tray"; \
	Filename: "{app}\TAMALabsTray.exe"; \
	WorkingDir: "{app}"

[Run]
; --- Register TAMALabs as Windows Service ---
Filename: "{app}\nssm.exe"; \
	Parameters: "install TAMALabs ""{app}\TAMALabs.exe"""; \
	Flags: runhidden waituntilterminated; \
	StatusMsg: "Registering TAMALabs service..."

Filename: "sc.exe"; \
	Parameters: "config TAMALabs start= auto"; \
	Flags: runhidden waituntilterminated; \
	StatusMsg: "Configuring TAMALabs service startup..."

Filename: "{app}\nssm.exe"; \
	Parameters: "start TAMALabs"; \
	Flags: runhidden waituntilterminated; \
	StatusMsg: "Starting TAMALabs service..."

; --- Configure Tray Auto-Start via Task Scheduler ---
Filename: "schtasks.exe"; \
	Parameters: "/create /tn ""TAMALabs Tray"" /tr ""\""{app}\TAMALabsTray.exe\"""" /sc onlogon /rl HIGHEST /f /ru ""%USERNAME%"""; \
	Flags: runhidden waituntilterminated; \
	StatusMsg: "Setting up TAMALabs tray auto-start with admin rights..."

; --- Start tray manually only if task not exists (avoid duplicate tray) ---
Filename: "{app}\TAMALabsTray.exe"; \
	Check: not TrayTaskExists; \
	Flags: nowait; \
	StatusMsg: "Starting TAMALabs system tray..."

[UninstallRun]
Filename: "taskkill.exe"; Parameters: "/F /IM TAMALabsTray.exe"; Flags: runhidden waituntilterminated; RunOnceId: "KillTray"
Filename: "{app}\nssm.exe"; Parameters: "stop TAMALabs"; Flags: runhidden waituntilterminated; RunOnceId: "StopService"
Filename: "{app}\nssm.exe"; Parameters: "remove TAMALabs confirm"; Flags: runhidden waituntilterminated; RunOnceId: "RemoveService"
Filename: "schtasks.exe"; Parameters: "/delete /tn ""TAMALabs Tray"" /f"; Flags: runhidden waituntilterminated; RunOnceId: "RemoveScheduledTask"
Filename: "cmd.exe"; Parameters: "/C rmdir /S /Q ""{app}"""; Flags: runhidden waituntilterminated; RunOnceId: "RemoveAppDir"

[UninstallDelete]
Type: files; Name: "{userstartup}\TAMALabs Tray.lnk"

[Code]
function TrayTaskExists: Boolean;
var
ResultCode: Integer;
begin
Exec('schtasks.exe', '/query /tn "TAMALabs Tray"', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
Result := (ResultCode = 0);
end;