import './style.css';

import {
  LoadProfiles, SaveProfile, DeleteProfile,
  StartTimer, ResumeTimer, PauseTimer, StopTimer, GetTimerState,
  PlayLooping, PlayShuffleFolder, StopAudio, SetVolume, GetAudioState,
  CheckResumeSession, PickMusicFile, PickMusicFolder,
  GetSettings, SaveSettings,
  GetStats, RecordSessionComplete
} from '../wailsjs/go/app/App';

import { EventsOn } from '../wailsjs/runtime/runtime';
import {
  WindowGetSize, WindowSetSize, WindowSetMinSize,
  WindowGetPosition, WindowSetPosition,
  WindowSetAlwaysOnTop
} from '../wailsjs/runtime/runtime';

// ── DOM refs ──────────────────────────────────────────────────────────────────
const timerEl        = document.getElementById('timerDisplay');
const fillEl         = document.getElementById('progressFill');
const startBtn       = document.getElementById('startBtn');
const stopBtn        = document.getElementById('stopBtn');
const skipBtn        = document.getElementById('skipBtn');
const profileSelect  = document.getElementById('profileSelect');
const resumeBanner   = document.getElementById('resumeBanner');
const resumeText     = document.getElementById('resumeText');
const resumeBtn      = document.getElementById('resumeBtn');
const audioDot       = document.getElementById('audioDot');
const trackNameEl    = document.getElementById('trackName');
const audioSubEl     = document.getElementById('audioSub');
const volumeSlider   = document.getElementById('volumeSlider');

// Profile panel
const overlay        = document.getElementById('overlay');
const profilePanel   = document.getElementById('profilePanel');
const profileList    = document.getElementById('profileList');
const profileForm    = document.getElementById('profileForm');
const formTitle      = document.getElementById('formTitle');
const pfName         = document.getElementById('pfName');
const pfDuration     = document.getElementById('pfDuration');
const pfMusicPath    = document.getElementById('pfMusicPath');
const pfShuffle      = document.getElementById('pfShuffle');
const pfEditId       = document.getElementById('pfEditId');
const pfBreakDuration  = document.getElementById('pfBreakDuration');
const pfBreakMusicPath = document.getElementById('pfBreakMusicPath');
const pfBreakShuffle   = document.getElementById('pfBreakShuffle');
const pfIsDefault      = document.getElementById('pfIsDefault');
const modeBadge        = document.getElementById('modeBadge');

// Mini widget
const miniWidget   = document.getElementById('miniWidget');
const miniTime     = document.getElementById('miniTime');
const miniBadge    = document.getElementById('miniBadge');
const miniPlayPause = document.getElementById('miniPlayPause');
const miniExpand   = document.getElementById('miniExpand');
const miniClose    = document.getElementById('miniClose');

// Settings panel
const settingsPanel  = document.getElementById('settingsPanel');
const stVolume       = document.getElementById('stVolume');
const stAutoAudio    = document.getElementById('stAutoAudio');
const stNotify       = document.getElementById('stNotify');
const stAutoNext     = document.getElementById('stAutoNext');
const stMinTray      = document.getElementById('stMinTray');
const stTheme        = document.getElementById('stTheme');
const settingsSaved  = document.getElementById('settingsSaved');

// ── App state ─────────────────────────────────────────────────────────────────
let profiles     = [];
let settings     = {};
let totalSec     = 25 * 60;
let remainSec    = totalSec;
let isRunning    = false;
let savedSession  = null;
let sessionType   = 'work'; // 'work' | 'break'
let activeProfile = null;   // currently running profile
let isMiniMode    = false;  // window is in compact mini-timer mode
let savedWindowState = null; // { width, height, x, y } before entering mini mode

// ── Helpers ───────────────────────────────────────────────────────────────────
function fmt(s) {
  const m   = Math.floor(s / 60).toString().padStart(2, '0');
  const sec = (s % 60).toString().padStart(2, '0');
  return `${m}:${sec}`;
}

function updateTimerUI(remaining, total) {
  timerEl.textContent = fmt(remaining);
  const pct = total > 0 ? ((total - remaining) / total) * 100 : 0;
  fillEl.style.width  = pct + '%';
}

function setRunningUI(running) {
  isRunning = running;
  startBtn.textContent       = running ? 'Pause' : 'Start';
  startBtn.style.background  = running ? 'rgba(255,180,50,0.55)' : '';
  miniPlayPause.textContent  = running ? '\u23F8' : '\u25B6';
}

// ── Panel helpers ─────────────────────────────────────────────────────────────
function openPanel(panel) {
  overlay.classList.add('show');
  panel.classList.add('open');
}
function closeAllPanels() {
  overlay.classList.remove('show');
  profilePanel.classList.remove('open');
  settingsPanel.classList.remove('open');
}

overlay.addEventListener('click', closeAllPanels);
document.getElementById('openProfiles').addEventListener('click', () => openPanel(profilePanel));
document.getElementById('openSettings').addEventListener('click', async () => {
  await loadSettingsIntoForm();
  openPanel(settingsPanel);
});
document.getElementById('closeProfiles').addEventListener('click', closeAllPanels);
document.getElementById('closeSettings').addEventListener('click', closeAllPanels);

// ── Profile list rendering ────────────────────────────────────────────────────
function renderProfileList() {
  profileList.innerHTML = profiles.map(p => `
    <div class="profile-item${p.isDefault ? ' is-default' : ''}">
      <div class="profile-item-info">
        <div class="profile-item-name">${p.isDefault ? '<span class="default-star" title="Default">★</span> ' : ''}${escHtml(p.name)}</div>
        <div class="profile-item-meta">${Math.floor(p.durationSec/60)} min${p.breakDurationSec > 0 ? ' + ' + Math.floor(p.breakDurationSec/60) + 'm break' : ''}${p.musicPath ? ' · ' + (p.shuffle ? 'Shuffle' : 'Loop') : ''}</div>
      </div>
      <div class="profile-item-actions">
        <button class="item-btn" data-id="${p.id}" data-action="edit">Edit</button>
        <button class="item-btn del" data-id="${p.id}" data-action="delete">Delete</button>
      </div>
    </div>
  `).join('');

  profileList.querySelectorAll('[data-action]').forEach(btn => {
    btn.addEventListener('click', () => {
      const id = btn.dataset.id;
      if (btn.dataset.action === 'edit') openEditForm(id);
      else confirmDeleteProfile(id);
    });
  });
}

function refreshDropdown(selectDefault = false) {
  const prev = profileSelect.value;
  profileSelect.innerHTML = profiles.map(p =>
    `<option value="${p.id}">${escHtml(p.name)}</option>`
  ).join('');
  if (selectDefault) {
    const def = profiles.find(p => p.isDefault);
    profileSelect.value = def ? def.id : (profiles[0]?.id || '');
  } else if (profileSelect.querySelector(`option[value="${prev}"]`)) {
    profileSelect.value = prev;
  }
  const sel = profiles.find(p => p.id === profileSelect.value) || profiles[0];
  if (sel) { totalSec = remainSec = sel.durationSec; updateTimerUI(remainSec, totalSec); }
}

function escHtml(s) {
  return (s || '').replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
}

// ── Profile form ──────────────────────────────────────────────────────────────
function showForm(show) {
  profileForm.style.display  = show ? 'flex' : 'none';
  document.getElementById('newProfileBtn').style.display = show ? 'none' : '';
}

function openNewForm() {
  formTitle.textContent      = 'New Profile';
  pfName.value               = '';
  pfDuration.value           = '25';
  pfMusicPath.value          = '';
  pfShuffle.checked          = false;
  pfBreakDuration.value      = '0';
  pfBreakMusicPath.value     = '';
  pfBreakShuffle.checked     = false;
  pfIsDefault.checked        = false;
  pfEditId.value             = '';
  showForm(true);
  pfName.focus();
}

function openEditForm(id) {
  const p = profiles.find(x => x.id === id);
  if (!p) return;
  formTitle.textContent      = 'Edit Profile';
  pfName.value               = p.name;
  pfDuration.value           = Math.floor(p.durationSec / 60).toString();
  pfMusicPath.value          = p.musicPath || '';
  pfShuffle.checked          = !!p.shuffle;
  pfBreakDuration.value      = Math.floor((p.breakDurationSec || 0) / 60).toString();
  pfBreakMusicPath.value     = p.breakMusicPath || '';
  pfBreakShuffle.checked     = !!p.breakShuffle;
  pfIsDefault.checked        = !!p.isDefault;
  pfEditId.value             = p.id;
  showForm(true);
}

document.getElementById('newProfileBtn').addEventListener('click', openNewForm);
document.getElementById('cancelProfileBtn').addEventListener('click', () => showForm(false));

document.getElementById('pickFile').addEventListener('click', async () => {
  const path = await PickMusicFile().catch(() => '');
  if (path) { pfMusicPath.value = path; pfShuffle.checked = false; }
});

document.getElementById('pickFolder').addEventListener('click', async () => {
  const path = await PickMusicFolder().catch(() => '');
  if (path) { pfMusicPath.value = path; pfShuffle.checked = true; }
});

document.getElementById('clearMusic').addEventListener('click', () => {
  pfMusicPath.value = '';
  pfShuffle.checked = false;
});

document.getElementById('pickBreakFile').addEventListener('click', async () => {
  const path = await PickMusicFile().catch(() => '');
  if (path) { pfBreakMusicPath.value = path; pfBreakShuffle.checked = false; }
});

document.getElementById('pickBreakFolder').addEventListener('click', async () => {
  const path = await PickMusicFolder().catch(() => '');
  if (path) { pfBreakMusicPath.value = path; pfBreakShuffle.checked = true; }
});

document.getElementById('clearBreakMusic').addEventListener('click', () => {
  pfBreakMusicPath.value = '';
  pfBreakShuffle.checked = false;
});

document.getElementById('saveProfileBtn').addEventListener('click', async () => {
  const name      = pfName.value.trim();
  if (!name) { pfName.focus(); return; }
  const dur       = Math.max(1, parseInt(pfDuration.value, 10) || 25);
  const breakMins = Math.max(0, parseInt(pfBreakDuration.value, 10) || 0);
  const id        = pfEditId.value || ('p' + Date.now());
  const p = {
    id,
    name,
    durationSec:      dur * 60,
    musicPath:        pfMusicPath.value.trim(),
    shuffle:          !!pfShuffle.checked,
    breakDurationSec: breakMins * 60,
    breakMusicPath:   pfBreakMusicPath.value.trim(),
    breakShuffle:     !!pfBreakShuffle.checked,
    isDefault:        !!pfIsDefault.checked,
  };
  await SaveProfile(p).catch(console.error);
  // If marked as default, clear isDefault on all others in local cache
  if (p.isDefault) profiles.forEach(x => { if (x.id !== id) x.isDefault = false; });
  // Update local cache
  const idx = profiles.findIndex(x => x.id === id);
  if (idx >= 0) profiles[idx] = p; else profiles.push(p);
  renderProfileList();
  refreshDropdown();
  showForm(false);
});

async function confirmDeleteProfile(id) {
  const p = profiles.find(x => x.id === id);
  if (!p) return;
  if (!confirm(`Delete profile "${p.name}"?`)) return;
  await DeleteProfile(id).catch(console.error);
  profiles = profiles.filter(x => x.id !== id);
  renderProfileList();
  refreshDropdown();
}

// ── Settings form ─────────────────────────────────────────────────────────────
async function loadSettingsIntoForm() {
  try {
    settings = await GetSettings();
    stVolume.value         = settings.defaultVolume ?? 70;
    stAutoAudio.checked    = !!settings.autoStartAudio;
    stNotify.checked       = !!settings.notifyOnComplete;
    stAutoNext.checked     = !!settings.autoStartNextTimer;
    stMinTray.checked      = !!settings.minimizeToTray;
    stTheme.value          = settings.theme || 'dark';
  } catch (e) { console.error('GetSettings failed', e); }
}

document.getElementById('saveSettingsBtn').addEventListener('click', async () => {
  const s = {
    defaultVolume:      parseInt(stVolume.value, 10),
    autoStartAudio:     stAutoAudio.checked,
    notifyOnComplete:   stNotify.checked,
    autoStartNextTimer: stAutoNext.checked,
    minimizeToTray:     stMinTray.checked,
    theme:              stTheme.value || 'dark',
  };
  await SaveSettings(s).catch(console.error);
  settings = s;
  applyTheme(s.theme);
  // Apply volume immediately
  SetVolume(s.defaultVolume).catch(() => {});
  volumeSlider.value = s.defaultVolume;
  // Show saved indicator
  settingsSaved.style.display = 'block';
  setTimeout(() => settingsSaved.style.display = 'none', 1800);
});

// ── Wails events ──────────────────────────────────────────────────────────────
EventsOn('timerTicked', (data) => {
  remainSec = data.remainingSec;
  updateTimerUI(remainSec, totalSec);
  // sync mini widget in both overlay and mini-mode
  if (isMiniMode || miniWidget.style.display !== 'none') {
    miniTime.textContent = fmt(remainSec);
  }
});

EventsOn('timerCompleted', async () => {
  setRunningUI(false);
  fillEl.style.width = '0%';

  if (sessionType === 'work') {
    // Work session done — record stats
    try { updateStatsUI(await RecordSessionComplete()); } catch (_) {}
    if (settings.notifyOnComplete) {
      try { new Notification('FocusPlay', { body: 'Work session complete! Take a break.' }); } catch (_) {}
    }
    // Auto-start break if profile has one configured
    if (activeProfile && activeProfile.breakDurationSec > 0) {
      setTimeout(() => startBreak(activeProfile), 800);
      return;
    }
  } else {
    // Break done — switch back to work mode
    sessionType = 'work';
    totalSec    = activeProfile ? activeProfile.durationSec : 25 * 60;
    updateModeBadge();
    if (settings.notifyOnComplete) {
      try { new Notification('FocusPlay', { body: "Break's over! Time to focus." }); } catch (_) {}
    }
  }

  remainSec = totalSec;
  updateTimerUI(remainSec, totalSec);
  if (settings.autoStartNextTimer) {
    const next = activeProfile || profiles.find(p => p.id === profileSelect.value) || profiles[0];
    if (next) setTimeout(() => startSession(next), 800);
  }
});

EventsOn('audioStateChanged', (data) => updateAudioUI(data));

function updateAudioUI(data) {
  if (!data) return;
  if (data.state === 'playing') {
    audioDot.classList.add('active');
    trackNameEl.textContent = data.trackName || 'Playing';
    audioSubEl.textContent  = data.trackInfo  || '';
  } else {
    audioDot.classList.remove('active');
    trackNameEl.textContent = data.trackName || 'No audio';
    audioSubEl.textContent  = data.state === 'stopped' ? 'Stopped' : '\u2014';
  }
}

function updateStatsUI(data) {
  if (!data) return;
  const countEl  = document.getElementById('sessionCount');
  const streakEl = document.getElementById('streakCount');
  if (countEl)  countEl.textContent  = data.sessionsToday ?? 0;
  if (streakEl) streakEl.textContent = (data.streak ?? 0) + (data.streak === 1 ? ' day' : ' days');
}

// ── Timer controls ────────────────────────────────────────────────────────────
async function startSession(profile) {
  activeProfile = profile;
  sessionType   = 'work';
  totalSec      = profile.durationSec;
  remainSec     = totalSec;
  updateTimerUI(remainSec, totalSec);
  setRunningUI(true);
  updateModeBadge();
  resumeBanner.style.display = 'none';
  await StartTimer(profile.id, profile.durationSec).catch(console.error);
  if (profile.musicPath && settings.autoStartAudio !== false) {
    if (profile.shuffle) await PlayShuffleFolder(profile.musicPath).catch(console.error);
    else                 await PlayLooping(profile.musicPath).catch(console.error);
  }
}

async function startBreak(profile) {
  sessionType = 'break';
  totalSec    = profile.breakDurationSec;
  remainSec   = totalSec;
  updateTimerUI(remainSec, totalSec);
  setRunningUI(true);
  updateModeBadge();
  await StartTimer(profile.id + '-break', profile.breakDurationSec).catch(console.error);
  if (settings.autoStartAudio !== false) {
    // Use break music if set, otherwise fall back to work music
    const musicPath = profile.breakMusicPath || profile.musicPath;
    const shuffle   = profile.breakMusicPath ? profile.breakShuffle : profile.shuffle;
    if (musicPath) {
      if (shuffle) await PlayShuffleFolder(musicPath).catch(console.error);
      else         await PlayLooping(musicPath).catch(console.error);
    } else {
      await StopAudio().catch(console.error);
    }
  }
}

function applyTheme(theme) {
  document.documentElement.setAttribute('data-theme', theme || 'dark');
}

function updateModeBadge() {
  if (!modeBadge) return;
  if (sessionType === 'break') {
    modeBadge.textContent = '\u25CF Break';
    modeBadge.className = 'badge is-break';
    miniBadge.textContent = '\u25CF Break';
    miniBadge.className = 'mini-badge is-break';
  } else {
    modeBadge.textContent = '\u25CF Work';
    modeBadge.className = 'badge';
    miniBadge.textContent = '\u25CF Work';
    miniBadge.className = 'mini-badge';
  }
}

// ── Mini widget (window-mode switching) ───────────────────────────────────────
const MINI_W = 220;
const MINI_H = 170;

async function showMini() {
  if (isMiniMode) return;
  // Save current window geometry
  const size = await WindowGetSize();
  const pos  = await WindowGetPosition();
  savedWindowState = { width: size.w, height: size.h, x: pos.x, y: pos.y };

  miniTime.textContent = fmt(remainSec);
  miniWidget.style.display = 'flex';
  document.body.classList.add('mini-mode');
  isMiniMode = true;

  WindowSetMinSize(MINI_W, MINI_H);
  WindowSetSize(MINI_W, MINI_H);
  WindowSetAlwaysOnTop(true);
}

async function hideMini() {
  if (!isMiniMode) return;
  document.body.classList.remove('mini-mode');
  miniWidget.style.display = 'none';
  isMiniMode = false;

  WindowSetAlwaysOnTop(false);
  if (savedWindowState) {
    WindowSetMinSize(420, 580);
    WindowSetSize(savedWindowState.width, savedWindowState.height);
    WindowSetPosition(savedWindowState.x, savedWindowState.y);
    savedWindowState = null;
  }
}

document.getElementById('openMini').addEventListener('click', () => {
  if (!isMiniMode) showMini();
  else hideMini();
});
miniExpand.addEventListener('click', hideMini);
miniClose.addEventListener('click', async () => {
  await hideMini();
});

miniPlayPause.addEventListener('click', () => {
  startBtn.click();
});

// Drag support for mini widget (only in overlay mode, not mini-mode)
(function() {
  const drag = document.getElementById('miniDrag');
  let dx = 0, dy = 0, startX = 0, startY = 0;
  drag.addEventListener('mousedown', (e) => {
    if (isMiniMode) return; // window-level drag in mini mode handled by --wails-draggable
    startX = e.clientX;
    startY = e.clientY;
    const rect = miniWidget.getBoundingClientRect();
    dx = e.clientX - rect.left;
    dy = e.clientY - rect.top;
    function onMove(ev) {
      const newLeft = ev.clientX - dx;
      const newTop  = ev.clientY - dy;
      miniWidget.style.left  = Math.max(0, Math.min(window.innerWidth  - miniWidget.offsetWidth,  newLeft)) + 'px';
      miniWidget.style.top   = Math.max(0, Math.min(window.innerHeight - miniWidget.offsetHeight, newTop))  + 'px';
      miniWidget.style.right = 'auto';
    }
    function onUp(ev) {
      document.removeEventListener('mousemove', onMove);
      document.removeEventListener('mouseup', onUp);
    }
    document.addEventListener('mousemove', onMove);
    document.addEventListener('mouseup', onUp);
  });
})();

startBtn.addEventListener('click', async () => {
  if (!isRunning) {
    const sel = profiles.find(p => p.id === profileSelect.value) || profiles[0];
    if (sel) await startSession(sel);
  } else {
    setRunningUI(false);
    await PauseTimer().catch(console.error);
    await StopAudio().catch(console.error);
  }
});

stopBtn.addEventListener('click', async () => {
  setRunningUI(false);
  sessionType = 'work';
  totalSec    = activeProfile ? activeProfile.durationSec : totalSec;
  remainSec   = totalSec;
  updateTimerUI(remainSec, totalSec);
  updateModeBadge();
  fillEl.style.width = '0%';
  await StopTimer().catch(console.error);
  await StopAudio().catch(console.error);
});

skipBtn.addEventListener('click', async () => {
  setRunningUI(false);
  await StopTimer().catch(console.error);

  if (sessionType === 'break') {
    // Skip break → back to work ready
    sessionType = 'work';
    totalSec    = activeProfile ? activeProfile.durationSec : 25 * 60;
    remainSec   = totalSec;
    updateTimerUI(remainSec, totalSec);
    fillEl.style.width = '0%';
    updateModeBadge();
    await StopAudio().catch(console.error);
  } else {
    // Skip work → if profile has a break, start it; else reset like stop
    fillEl.style.width = '0%';
    if (activeProfile && activeProfile.breakDurationSec > 0) {
      setTimeout(() => startBreak(activeProfile), 300);
    } else {
      totalSec  = activeProfile ? activeProfile.durationSec : 25 * 60;
      remainSec = totalSec;
      updateTimerUI(remainSec, totalSec);
      await StopAudio().catch(console.error);
    }
  }
});

resumeBtn.addEventListener('click', async () => {
  if (!savedSession) return;
  resumeBanner.style.display = 'none';
  totalSec  = savedSession.totalSec;
  remainSec = savedSession.remainingSec;
  updateTimerUI(remainSec, totalSec);
  setRunningUI(true);
  await ResumeTimer(savedSession).catch(console.error);
  // Restart audio for the resumed profile
  const prof = profiles.find(p => p.id === savedSession.profileId);
  if (prof && prof.musicPath && settings.autoStartAudio !== false) {
    if (prof.shuffle) await PlayShuffleFolder(prof.musicPath).catch(console.error);
    else              await PlayLooping(prof.musicPath).catch(console.error);
  }
});

profileSelect.addEventListener('change', async () => {
  const sel = profiles.find(p => p.id === profileSelect.value);
  if (!sel) return;
  if (isRunning) { await StopTimer().catch(console.error); setRunningUI(false); }
  // Always stop audio when switching profiles
  await StopAudio().catch(console.error);
  totalSec  = sel.durationSec;
  remainSec = totalSec;
  updateTimerUI(remainSec, totalSec);
  fillEl.style.width = '0%';
});

volumeSlider.addEventListener('input', async () => {
  await SetVolume(parseInt(volumeSlider.value, 10)).catch(console.error);
});

// ── Keyboard shortcuts ────────────────────────────────────────────────────────
document.addEventListener('keydown', (e) => {
  // Ignore when typing in inputs
  const tag = (e.target.tagName || '').toLowerCase();
  if (tag === 'input' || tag === 'textarea' || tag === 'select') return;
  // Ignore when a panel is open
  if (overlay.classList.contains('show')) return;

  switch (e.code) {
    case 'Space':
      e.preventDefault();
      startBtn.click();
      break;
    case 'Escape':
      e.preventDefault();
      stopBtn.click();
      break;
    case 'KeyS':
      e.preventDefault();
      skipBtn.click();
      break;
    case 'KeyM':
      e.preventDefault();
      document.getElementById('openMini').click();
      break;
  }
});

// ── Boot ──────────────────────────────────────────────────────────────────────
async function init() {
  // Request notification permission early (WebView2 may silently deny otherwise)
  if ('Notification' in window && Notification.permission === 'default') {
    Notification.requestPermission().catch(() => {});
  }

  // Load settings first
  try {
    settings = await GetSettings();
    volumeSlider.value = settings.defaultVolume ?? 70;
    SetVolume(settings.defaultVolume ?? 70).catch(() => {});
    applyTheme(settings.theme);
  } catch (e) {}

  // Load profiles
  try {
    profiles = await LoadProfiles();
    refreshDropdown(true);
    renderProfileList();
  } catch (e) { console.error('LoadProfiles failed', e); }

  // Check resume session
  try {
    savedSession = await CheckResumeSession();
    if (savedSession && savedSession.remainingSec > 0) {
      resumeText.textContent     = `Previous session found \u2014 ${fmt(savedSession.remainingSec)} remaining`;
      resumeBanner.style.display = 'flex';
      totalSec  = savedSession.totalSec;
      remainSec = savedSession.remainingSec;
      updateTimerUI(remainSec, totalSec);
    }
  } catch (e) {}

  // Sync live timer state
  try {
    const state = await GetTimerState();
    if (state.running) {
      totalSec  = state.totalSec;
      remainSec = state.remainingSec;
      updateTimerUI(remainSec, totalSec);
      setRunningUI(true);
    }
  } catch (e) {}

  // Sync audio
  try { updateAudioUI(await GetAudioState()); } catch (e) {}

  // Load footer stats
  try { updateStatsUI(await GetStats()); } catch (e) {}
}

init();
