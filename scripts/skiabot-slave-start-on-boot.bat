if not defined SKIA_REPO_DIR (set SKIA_REPO_DIR=%HOMEDRIVE%HOMEPATH%)
chdir /d %SKIA_REPO_DIR%
call gclient config --spec="solutions = [{ 'name': 'buildbot','url': 'https://skia.googlesource.com/buildbot.git','deps_file': 'DEPS','managed': True,'custom_deps': {},'safesync_url': '',},{ 'name': 'src','url': 'https://chromium.googlesource.com/chromium/src.git','deps_file': '.DEPS.git','managed': True,'custom_deps': {},'safesync_url': '',},]"
call gclient sync --force -j1
cd buildbot
python scripts\launch_slaves.py