import './style.css';

import {
  LoadProfiles, SaveProfile, DeleteProfile,
  StartTimer, ResumeTimer, PauseTimer, StopTimer, GetTimerState,
  PlayLooping, PlayShuffleFolder, StopAudio, SetVolume, GetAudioState,
  CheckResumeSession, PickMusicFile, PickMusicFolder,
  GetSettings, SaveSettings,
  GetStats, RecordSessionComplete
} from '../wailsjs/go/main/App';

import { EventsOn } from '../wailsjs/runtime/runtime';

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

// Settings panel
const settingsPanel  = document.getElementById('settingsPanel');
const stVolume       = document.getElementById('stVolume');
const stAutoAudio    = document.getElementById('stAutoAudio');
const stNotify       = document.getElementById('stNotify');
const stAutoNext     = document.getElementById('stAutoNext');
const stMinTray      = document.getElementById('stMinTray');
const settingsSaved  = document.getElementById('settingsSaved');

// ── App state ─────────────────────────────────────────────────────────────────
let profiles     = [];
let settings     = {};
let totalSec     = 25 * 60;
let remainSec    = totalSec;
let isRunning    = false;
let savedSession = null;

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
    <div class="profile-item">
      <div class="profile-item-info">
        <div class="profile-item-name">${escHtml(p.name)}</div>
        <div class="profile-item-meta">${Math.floor(p.durationSec/60)} min${p.musicPath ? ' · ' + (p.shuffle ? 'Shuffle' : 'Loop') : ''}</div>
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

function refreshDropdown() {
  const prev = profileSelect.value;
  profileSelect.innerHTML = profiles.map(p =>
    `<option value="${p.id}">${escHtml(p.name)}</option>`
  ).join('');
  if (profileSelect.querySelector(`option[value="${prev}"]`)) profileSelect.value = prev;
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
  formTitle.textContent = 'New Profile';
  pfName.value      = '';
  pfDuration.value  = '25';
  pfMusicPath.value = '';
  pfShuffle.checked = false;
  pfEditId.value    = '';
  showForm(true);
  pfName.focus();
}

function openEditForm(id) {
  const p = profiles.find(x => x.id === id);
  if (!p) return;
  formTitle.textContent = 'Edit Profile';
  pfName.value      = p.name;
  pfDuration.value  = Math.floor(p.durationSec / 60).toString();
  pfMusicPath.value = p.musicPath || '';
  pfShuffle.checked = !!p.shuffle;
  pfEditId.value    = p.id;
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

document.getElementById('saveProfileBtn').addEventListener('click', async () => {
  const name = pfName.value.trim();
  if (!name) { pfName.focus(); return; }
  const dur  = Math.max(1, parseInt(pfDuration.value, 10) || 25);
  const id   = pfEditId.value || ('p' + Date.now());
  const p = {
    id,
    name,
    durationSec: dur * 60,
    musicPath:   pfMusicPath.value.trim(),
    shuffle:     !!pfShuffle.checked
  };
  await SaveProfile(p).catch(console.error);
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
  } catch (e) { console.error('GetSettings failed', e); }
}

document.getElementById('saveSettingsBtn').addEventListener('click', async () => {
  const s = {
    defaultVolume:      parseInt(stVolume.value, 10),
    autoStartAudio:     stAutoAudio.checked,
    notifyOnComplete:   stNotify.checked,
    autoStartNextTimer: stAutoNext.checked,
    minimizeToTray:     stMinTray.checked
  };
  await SaveSettings(s).catch(console.error);
  settings = s;
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
});

EventsOn('timerCompleted', async () => {
  setRunningUI(false);
  remainSec = totalSec;
  updateTimerUI(remainSec, totalSec);
  fillEl.style.width = '0%';
  // Record the completed session and refresh footer stats
  try { updateStatsUI(await RecordSessionComplete()); } catch (_) {}
  if (settings.notifyOnComplete) {
    try { new Notification('FocusPlay', { body: 'Session complete!' }); } catch (_) {}
  }
  if (settings.autoStartNextTimer) {
    const sel = profiles.find(p => p.id === profileSelect.value) || profiles[0];
    if (sel) startSession(sel);
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
  totalSec  = profile.durationSec;
  remainSec = totalSec;
  updateTimerUI(remainSec, totalSec);
  setRunningUI(true);
  resumeBanner.style.display = 'none';
  await StartTimer(profile.id, profile.durationSec).catch(console.error);
  if (profile.musicPath && settings.autoStartAudio !== false) {
    if (profile.shuffle) await PlayShuffleFolder(profile.musicPath).catch(console.error);
    else                 await PlayLooping(profile.musicPath).catch(console.error);
  }
}

startBtn.addEventListener('click', async () => {
  if (!isRunning) {
    const sel = profiles.find(p => p.id === profileSelect.value) || profiles[0];
    if (sel) await startSession(sel);
  } else {
    setRunningUI(false);
    await PauseTimer().catch(console.error);
  }
});

stopBtn.addEventListener('click', async () => {
  setRunningUI(false);
  remainSec = totalSec;
  updateTimerUI(remainSec, totalSec);
  fillEl.style.width = '0%';
  await StopTimer().catch(console.error);
  await StopAudio().catch(console.error);
});

skipBtn.addEventListener('click', async () => {
  setRunningUI(false);
  remainSec = 0;
  updateTimerUI(0, totalSec);
  fillEl.style.width = '100%';
  await StopTimer().catch(console.error);
});

resumeBtn.addEventListener('click', async () => {
  if (!savedSession) return;
  resumeBanner.style.display = 'none';
  totalSec  = savedSession.totalSec;
  remainSec = savedSession.remainingSec;
  updateTimerUI(remainSec, totalSec);
  setRunningUI(true);
  await ResumeTimer(savedSession).catch(console.error);
});

profileSelect.addEventListener('change', async () => {
  const sel = profiles.find(p => p.id === profileSelect.value);
  if (!sel) return;
  if (isRunning) { await StopTimer().catch(console.error); await StopAudio().catch(console.error); setRunningUI(false); }
  totalSec  = sel.durationSec;
  remainSec = totalSec;
  updateTimerUI(remainSec, totalSec);
  fillEl.style.width = '0%';
});

volumeSlider.addEventListener('input', async () => {
  await SetVolume(parseInt(volumeSlider.value, 10)).catch(console.error);
});

// ── Boot ──────────────────────────────────────────────────────────────────────
async function init() {
  // Load settings first
  try {
    settings = await GetSettings();
    volumeSlider.value = settings.defaultVolume ?? 70;
    SetVolume(settings.defaultVolume ?? 70).catch(() => {});
  } catch (e) {}

  // Load profiles
  try {
    profiles = await LoadProfiles();
    refreshDropdown();
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
