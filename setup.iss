; ==================================================================
; Inno Setup Script for TAMALabs (Normal App + Auto-Start Tray)
; ==================================================================

; Read version from version.go
#define AppVer ReadIni(SourcePath + "version.ini", "Version", "AppVersion", "v1.0.0")

[Setup]
AppId={{A8A93F44-8D2B-4D75-9CC8-0C52B2184AC4}}
AppName=TAMALabs
AppVersion={#AppVer}
AppPublisher=Elga Tama
AppPublisherURL=https://tamalabs.biosystems.id/
AppSupportURL=https://tamalabs.biosystems.id/support
AppUpdatesURL=https://tamalabs.biosystems.id/updates
DefaultDirName={autopf}\TAMALabs
DefaultGroupName=TAMALabs
AllowNoIcons=yes
OutputBaseFilename=TAMALabs-setup-{#AppVer}
OutputDir=.\installers
Compression=lzma
SolidCompression=yes
WizardStyle=modern
ArchitecturesInstallIn64BitMode=x64compatible
PrivilegesRequired=admin
PrivilegesRequiredOverridesAllowed=dialog

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Code]
var
  LabInfoPage: TWizardPage;
  HospitalNameEdit: TEdit;
  LocationEdit: TEdit;

procedure InitializeWizard;
begin
  { Create custom page for Lab Information }
  LabInfoPage := CreateCustomPage(wpSelectDir, 
    'Laboratory Information', 
    'Please enter your laboratory information for integration service');
  
  { Hospital Name }
  with TLabel.Create(LabInfoPage) do
  begin
    Parent := LabInfoPage.Surface;
    Caption := 'Hospital Name:';
    Left := 0;
    Top := 0;
  end;
  
  HospitalNameEdit := TEdit.Create(LabInfoPage);
  with HospitalNameEdit do
  begin
    Parent := LabInfoPage.Surface;
    Left := 0;
    Top := 20;
    Width := 400;
    Text := '';
  end;
  
  { Location }
  with TLabel.Create(LabInfoPage) do
  begin
    Parent := LabInfoPage.Surface;
    Caption := 'Location:';
    Left := 0;
    Top := 60;
  end;
  
  LocationEdit := TEdit.Create(LabInfoPage);
  with LocationEdit do
  begin
    Parent := LabInfoPage.Surface;
    Left := 0;
    Top := 80;
    Width := 400;
    Text := '';
  end;
end;

function NextButtonClick(CurPageID: Integer): Boolean;
begin
  Result := True;
  
  if CurPageID = LabInfoPage.ID then
  begin
    { Validate inputs }
    if Trim(HospitalNameEdit.Text) = '' then
    begin
      MsgBox('Hospital Name cannot be empty!', mbError, MB_OK);
      Result := False;
      Exit;
    end;
    
    if Trim(LocationEdit.Text) = '' then
    begin
      MsgBox('Location cannot be empty!', mbError, MB_OK);
      Result := False;
      Exit;
    end;
  end;
end;

procedure CurStepChanged(CurStep: TSetupStep);
var
  LabInfoFile: String;
  LabInfoContent: TStringList;
begin
  if CurStep = ssPostInstall then
  begin
    { Create lab_info.json file }
    LabInfoFile := ExpandConstant('{app}\integration-service\lab_info.json');
    LabInfoContent := TStringList.Create;
    try
      LabInfoContent.Add('{');
      LabInfoContent.Add('  "lab_id": "",');
      LabInfoContent.Add('  "hospital_name": "' + HospitalNameEdit.Text + '",');
      LabInfoContent.Add('  "location": "' + LocationEdit.Text + '"');
      LabInfoContent.Add('}');
      LabInfoContent.SaveToFile(LabInfoFile);
    finally
      LabInfoContent.Free;
    end;
  end;
end;

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"

[Dirs]
Name: "{app}\logs"; Permissions: users-modify
Name: "{app}\integration-service\logs"; Permissions: users-modify

[Files]
Source: "bin\TAMALabs.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: ".env"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\TAMALabsTray.exe"; DestDir: "{app}"; Flags: ignoreversion

; Copy integration service executable
Source: "bin\TAMALabsIntegration.exe"; DestDir: "{app}\integration-service"; Flags: ignoreversion
Source: "integration-service\.env"; DestDir: "{app}\integration-service"; Flags: ignoreversion skipifsourcedoesntexist

[Icons]
; Optional desktop shortcut
Name: "{autodesktop}\TAMALabs (Open Web)"; \
    Filename: "rundll32.exe"; \
    Parameters: "url.dll,FileProtocolHandler http://127.0.0.1:8322"; \
    IconFilename: "{app}\TAMALabs.exe"; \
    Tasks: desktopicon

; Auto-start main app via Startup folder (for all users)
Name: "{commonstartup}\TAMALabs"; \
    Filename: "{app}\TAMALabs.exe"; \
    WorkingDir: "{app}"; \
    IconFilename: "{app}\TAMALabs.exe"

; Auto-start tray via Startup folder (for all users)
Name: "{commonstartup}\TAMALabs Tray"; \
    Filename: "{app}\TAMALabsTray.exe"; \
    WorkingDir: "{app}"; \
    IconFilename: "{app}\TAMALabsTray.exe"

; Auto-start integration service via Startup folder (for all users)
Name: "{commonstartup}\TAMALabs Integration Service"; \
    Filename: "{app}\integration-service\TAMALabsIntegration.exe"; \
    WorkingDir: "{app}\integration-service"; \
    IconFilename: "{app}\integration-service\TAMALabsIntegration.exe"

[Run]
; Jalankan aplikasi utama langsung (tanpa muncul checkbox di akhir)
Filename: "{app}\TAMALabs.exe"; \
    Description: "Start TAMALabs"; \
    Flags: nowait skipifsilent; \
    StatusMsg: "Starting TAMALabs..."

; Jalankan tray langsung (tanpa muncul checkbox di akhir)
Filename: "{app}\TAMALabsTray.exe"; \
    Description: "Start TAMALabs Tray"; \
    Flags: nowait skipifsilent; \
    WorkingDir: "{app}"; \
    StatusMsg: "Starting TAMALabs Tray..."

; Jalankan integration service langsung
Filename: "{app}\integration-service\TAMALabsIntegration.exe"; \
    Description: "Start TAMALabs Integration Service"; \
    Flags: nowait skipifsilent; \
    WorkingDir: "{app}\integration-service"; \
    StatusMsg: "Starting TAMALabs Integration Service..."

[UninstallRun]
; Tutup tray dan app sebelum uninstall
Filename: "taskkill.exe"; Parameters: "/F /IM TAMALabs.exe /T"; Flags: runhidden waituntilterminated; RunOnceId: "KillMain"
Filename: "taskkill.exe"; Parameters: "/F /IM TAMALabsTray.exe /T"; Flags: runhidden waituntilterminated; RunOnceId: "KillTray"
Filename: "taskkill.exe"; Parameters: "/F /IM TAMALabsIntegration.exe /T"; Flags: runhidden waituntilterminated; RunOnceId: "KillIntegration"

; Hapus folder aplikasi
Filename: "cmd.exe"; Parameters: "/C rmdir /S /Q ""{app}"""; Flags: runhidden waituntilterminated; RunOnceId: "RemoveAppFolder"

[UninstallDelete]
; Hapus shortcut startup & desktop
Type: files; Name: "{commonstartup}\TAMALabs Tray.lnk"
Type: files; Name: "{commonstartup}\TAMALabs Integration Service.lnk"
Type: files; Name: "{autodesktop}\TAMALabs (Open Web).lnk"
Type: dirifempty; Name: "{app}\logs"
Type: dirifempty; Name: "{app}\integration-service\logs"
Type: dirifempty; Name: "{app}\integration-service\tmp"
