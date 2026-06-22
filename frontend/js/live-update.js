// ============================================================================
// CONFIGURATION
// ============================================================================
const API_URL = 'http://localhost:8080/pl';
const INTERVAL_ONE_HOUR = 60 * 60 * 1000;
const RECONNECT_DELAY = 5000; // 5 Sekunden Wartezeit für Reconnects

// ============================================================================
// WEBSOCKET 1: SPECTATOR SERVER (Port 8080)
// ============================================================================
function connectSpectatorServer() {
    const socket = new WebSocket('ws://localhost:8080/ws');

    socket.onopen = () => {
        console.log('WebSocket-Connection to the Spectator Server established.');
    };

    socket.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            updateGameMeta(data);
            updateScoreDisplay(data);
            updateTimers(data.timers);
            updatePlayerScoreboard(data);
        } catch (error) {
            console.error('Error parsing the WebSocket data:', error);
        }
    };

    socket.onclose = () => {
        console.warn(`Spectator Server connection closed. Retrying in ${RECONNECT_DELAY / 1000}s...`);
        setTimeout(connectSpectatorServer, RECONNECT_DELAY); // Verbindet sich selbst neu, OHNE die Seite neu zu laden
    };

    socket.onerror = (error) => {
        console.error('WebSocket error (Spectator Server):', error);
    };
}

// ============================================================================
// WEBSOCKET 2: PYTHON SCRIPT (Port 8765)
// ============================================================================
function connectPythonScript() {
    const scriptSocket = new WebSocket('ws://localhost:8765/ws');

    scriptSocket.onopen = () => {
        console.log('WebSocket-Connection for the P-Script to the Spectator Server established.');
    };

    scriptSocket.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            updateTeamGold(data);
            updatePlayerGoldDiff(data);
            updatePlayerCs(data);
            updatePlayerBounty(data);
        } catch (error) {
            console.error('Error parsing the WebSocket for P-Script data:', error);
        }
    };

    scriptSocket.onclose = () => {
        console.warn(`P-Script connection closed. Retrying in ${RECONNECT_DELAY / 1000}s...`);
        setTimeout(connectPythonScript, RECONNECT_DELAY); // Verbindet sich selbst neu, OHNE die Seite neu zu laden
    };

    scriptSocket.onerror = (error) => {
        console.error('WebSocket for P-Script error:', error);
    };
}

// ============================================================================
// PRIME LEAGUE API (HTTP GET / Async-Interval)
// ============================================================================
async function fetchAndUpdateStandings() {
    try {
        const response = await fetch(API_URL);
        if (!response.ok) {
            throw new Error(`API-Fehler: Status ${response.status}`);
        }

        const data = await response.json();
        console.log("Empfangene Prime League Daten:", data);
        updateStandingsTable(data);
        updateScoreTeams(data);
    } catch (error) {
        // Schlägt die API fehl, bleibt die Tabelle einfach unverändert und das Script läuft weiter
        console.error("Error loading Prime League data (Will retry in 1 hour): ", error);
    }
}

// ============================================================================
// FORMATTING & HELPERS
// ============================================================================
function formatGameTime(seconds) {
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
}

function formatGold(goldAmount) {
    return `${(goldAmount / 1000).toFixed(1)}K`;
}

// ============================================================================
// DOM RENDERING FUNCTIONS
// ============================================================================
function updateGameMeta(data) {
    if (!data || !data.gameTime) return;
    const timeEl = document.querySelector('.game-time');
    if (timeEl) {
        timeEl.textContent = formatGameTime(data.gameTime);
    }
}

function updateScoreDisplay(data) {
    const blue = data.blueTeam;
    const red = data.redTeam;
    if (!blue || !red) return;

    document.querySelector('.blue-kills').textContent = blue.score;
    document.querySelector('.red-kills').textContent = red.score;

    if (blue.objectives) {
        blue.objectives.forEach(obj => {
            if (obj.key === 'turrets') document.querySelector('.blue-towers .obj-count').textContent = obj.kills;
            if (obj.key === 'grubs') document.querySelector('.blue-grubs .obj-count').textContent = obj.kills;
            if (obj.key === 'dragon') {
                const blueDrakesContainer = document.querySelector('.blue-drakes');
                if (blueDrakesContainer) {
                    blueDrakesContainer.innerHTML = `<img src="${obj.icon}" class="drake-icon" alt="${obj.name}">`;
                }
            }
        });
    }

    if (red.objectives) {
        red.objectives.forEach(obj => {
            if (obj.key === 'turrets') document.querySelector('.red-towers .obj-count').textContent = obj.kills;
            if (obj.key === 'grubs') document.querySelector('.red-grubs .obj-count').textContent = obj.kills;
            if (obj.key === 'dragon') {
                const redDrakesContainer = document.querySelector('.red-drakes');
                if (redDrakesContainer) {
                    redDrakesContainer.innerHTML = `<img src="${obj.icon}" class="drake-icon" alt="${obj.name}">`;
                }
            }
        });
    }
}

function updateTimers(timersArray) {
    if (!timersArray || !Array.isArray(timersArray)) return;

    timersArray.forEach(timer => {
        if (timer.type === 'dragon') {
            const drakeTimeEl = document.querySelector('#drake-timer-box .spawn-time');
            if (drakeTimeEl) {
                drakeTimeEl.textContent = timer.alive ? "LIVE" : formatGameTime(timer.SpawnTime);
                if (timer.alive) {
                    drakeTimeEl.classList.add('live');
                } else {
                    drakeTimeEl.classList.remove('live');
                }
            }
        }

        if (timer.type === 'baron') {
            const baronTimeEl = document.querySelector('#void-timer-box .spawn-time');
            if (baronTimeEl) {
                baronTimeEl.textContent = timer.alive ? "LIVE" : formatGameTime(timer.SpawnTime);
                if (timer.alive) {
                    baronTimeEl.classList.add('live');
                } else {
                    baronTimeEl.classList.remove('live');
                }
            }
        }
    });
}

function updatePlayerScoreboard(data) {
    const bluePlayers = data.blueTeam.players;
    const redPlayers = data.redTeam.players;
    if (!bluePlayers || !redPlayers) return;

    const rows = document.querySelectorAll('#player-scoreboard .player-row');

    for (let i = 0; i < 5; i++) {
        const row = rows[i];
        if (!row) break;

        const bluePlayerData = bluePlayers[i];
        const redPlayerData = redPlayers[i];

        const blueCard = row.querySelector('.blue-player');
        if (blueCard && bluePlayerData) updatePlayerCard(blueCard, bluePlayerData);

        const redCard = row.querySelector('.red-player');
        if (redCard && redPlayerData) updatePlayerCard(redCard, redPlayerData);
    }
}

function updatePlayerCard(cardEl, p) {
    const champImg = cardEl.querySelector('.champ-img');
    if (champImg) champImg.src = `${p.icon}`;

    const lvlEl = cardEl.querySelector('.champ-level');
    if (lvlEl) lvlEl.textContent = p.level;

    const nameEl = cardEl.querySelector('.player-name');
    if (nameEl) nameEl.textContent = p.riotId ? p.riotId.split('#')[0] : '';

    if (p.scores) {
        const kdaEl = cardEl.querySelector('.player-kda');
        if (kdaEl) kdaEl.textContent = `${p.scores.kills}/${p.scores.deaths}/${p.scores.assists}`;
    }

    const avatarBox = cardEl.querySelector('.champ-avatar-box');
    if (avatarBox) {
        const deathOverlay = cardEl.querySelector('.death-timer-overlay');
        const hpBar = cardEl.querySelector('.hp-bar');
        const manaBar = cardEl.querySelector('.mana-bar');

        if (p.isDead) {
            avatarBox.classList.add('dead');
            if (hpBar) {
                hpBar.style.width = '0%';
                hpBar.classList.add('dead-bar');
            }
            if (manaBar) manaBar.style.width = '0%';
            if (deathOverlay) {
                deathOverlay.classList.remove('hidden');
                deathOverlay.innerHTML = `<span>${Math.floor(p.respawnTimer)}</span>`;
            }
        } else {
            avatarBox.classList.remove('dead');
            if (hpBar) {
                hpBar.style.width = '100%';
                hpBar.classList.remove('dead-bar');
            }
            if (manaBar) manaBar.style.width = '100%';
            if (deathOverlay) deathOverlay.classList.add('hidden');
        }
    }

    if (p.runes) {
        const runeImages = cardEl.querySelectorAll('.runes-box img');
        if (runeImages.length >= 2) {
            if (p.runes.keystone) runeImages[0].src = p.runes.keystone.icon;
            if (p.runes.secondary) runeImages[1].src = p.runes.secondary.icon;
        }
    }

    const summSlots = cardEl.querySelectorAll('.summs-container .summ-slot');
    if (summSlots.length >= 2 && p.spells && p.spells.length >= 2) {
        summSlots[0].innerHTML = `<img src="${p.spells[0].icon}" alt="${p.spells[0].displayName}">`;
        summSlots[1].innerHTML = `<img src="${p.spells[1].icon}" alt="${p.spells[1].displayName}">`;
    }

    const itemSlots = cardEl.querySelectorAll('.items-line .item-slot');
    itemSlots.forEach(slot => {
        slot.className = 'item-slot empty';
        slot.innerHTML = '';
    });

    if (p.items && Array.isArray(p.items)) {
        p.items.forEach(item => {
            const slotIdx = item.slot;
            if (slotIdx >= 6) return;
            const slot = itemSlots[slotIdx];
            if (!slot) return;

            slot.className = 'item-slot';
            slot.innerHTML = `<img src="${item.icon}" alt="${item.id}">`;
            if (item.consumable) {
                slot.innerHTML += `<span class="item-count">${item.count}</span>`;
            }
        });
    }

    const wardSlot = cardEl.querySelector('.ward-slot');
    if (wardSlot && p.items) {
        const wardItem = p.items.find(item => item.slot === 6);
        if (wardItem) {
            wardSlot.classList.remove("empty-slot");
            wardSlot.innerHTML = `
                <img src="${wardItem.icon}" alt="${wardItem.id}">
                <span class="ward-count">${Math.floor(p.scores.wardScore)}</span>
            `;
        } else {
            wardSlot.innerHTML = `<span class="ward-count">${Math.floor(p.scores.wardScore)}</span>`;
        }
    }

    const ultSlot = cardEl.querySelector('.ult-slot');
    if (ultSlot && p.ultIcon) {
        ultSlot.innerHTML = `<img src="${p.ultIcon}" alt="U">`;
        if (p.level >= 6) {
            ultSlot.classList.remove("unlearned");
            ultSlot.classList.add("ready");
        } else {
            ultSlot.classList.remove("ready");
            ultSlot.classList.add("unlearned");
        }
    }
}

function updateTeamGold(data) {
    if (!data || !data.team_gold) return;
    const blue_gold = data.team_gold.blue_team;
    const red_gold = data.team_gold.red_team;

    document.querySelector('.blue-side .team-gold').textContent = blue_gold ? formatGold(blue_gold) : "0.0K";
    document.querySelector('.red-side .team-gold').textContent = red_gold ? formatGold(red_gold) : "0.0K";

    const goldDiff = blue_gold - red_gold;
    const blueLeadEl = document.querySelector('.blue-side .gold-lead');
    const redLeadEl = document.querySelector('.red-side .gold-lead');

    if (goldDiff > 0) {
        if (blueLeadEl) {
            blueLeadEl.classList.remove('hidden');
            blueLeadEl.classList.add('blue-side-lead');
            blueLeadEl.textContent = `+${formatGold(goldDiff)}`;
        }
        if (redLeadEl) redLeadEl.classList.add('hidden');
    } else if (goldDiff < 0) {
        if (redLeadEl) {
            redLeadEl.classList.remove('hidden');
            redLeadEl.classList.add('red-side-lead');
            redLeadEl.textContent = `+${formatGold(Math.abs(goldDiff))}`;
        }
        if (blueLeadEl) blueLeadEl.classList.add('hidden');
    } else {
        if (blueLeadEl) blueLeadEl.classList.add('hidden');
        if (redLeadEl) redLeadEl.classList.add('hidden');
    }
}

function updatePlayerGoldDiff(data) {
    if (!data || !data.gold || !data.gold.blue || !data.gold.red) return;
    const bluePlayers = data.gold.blue;
    const redPlayers = data.gold.red;

    const rows = document.querySelectorAll('#player-scoreboard .player-row');

    for (let i = 0; i < 5; i++) {
        const row = rows[i];
        if (!row) break;

        const bPlayer = `b${i+1}`;
        const rPlayer = `r${i+1}`;
        const bluePlayerGold = parseInt(bluePlayers[bPlayer], 10) || 0;
        const redPlayerGold = parseInt(redPlayers[rPlayer], 10) || 0;

        const goldDiff = bluePlayerGold - redPlayerGold;
        const LeadEl = row.querySelector('.lane-gold-diff');

        if (LeadEl) {
            const amountEl = LeadEl.querySelector('.diff-amount');
            LeadEl.classList.remove('blue-lead', 'red-lead', 'no-lead');

            if (goldDiff > 0) {
                LeadEl.classList.add('blue-lead');
                if (amountEl) amountEl.textContent = `+${formatGold(goldDiff)}`;
            } else if (goldDiff < 0) {
                LeadEl.classList.add('red-lead');
                if (amountEl) amountEl.textContent = `+${formatGold(Math.abs(goldDiff))}`;
            } else {
                LeadEl.classList.add('no-lead');
                if (amountEl) amountEl.textContent = '0.0K';
            }
        }
    }
}

function updatePlayerCs(data) {
    if (!data || !data.creep || !data.creep.blue || !data.creep.red) return;
    const bluePlayers = data.creep.blue;
    const redPlayers = data.creep.red;

    const rows = document.querySelectorAll('#player-scoreboard .player-row');

    for (let i = 0; i < 5; i++) {
        const row = rows[i];
        if (!row) break;

        const bPlayer = `b${i+1}`;
        const rPlayer = `r${i+1}`;

        const blueCard = row.querySelector('.blue-player');
        if (blueCard && bluePlayers[bPlayer] !== undefined) {
            const csEl = blueCard.querySelector('.player-cs');
            if (csEl) csEl.textContent = bluePlayers[bPlayer];
        }

        const redCard = row.querySelector('.red-player');
        if (redCard && redPlayers[rPlayer] !== undefined) {
            const csEl = redCard.querySelector('.player-cs');
            if (csEl) csEl.textContent = redPlayers[rPlayer];
        }
    }
}

function updatePlayerBounty(data) {
    if (!data || !data.bounty || !data.bounty.blue || !data.bounty.red) return;
    const bluePlayers = data.bounty.blue;
    const redPlayers = data.bounty.red;

    const rows = document.querySelectorAll('#player-scoreboard .player-row');

    for (let i = 0; i < 5; i++) {
        const row = rows[i];
        if (!row) break;

        const bPlayer = `b${i+1}`;
        const rPlayer = `r${i+1}`;

        const blueCard = row.querySelector('.blue-player');
        if (blueCard) {
            const bountyEl = blueCard.querySelector('.bounty-tag');
            if (bountyEl) {
                const bounty = bluePlayers[bPlayer] || 0;
                if (bounty > 0) {
                    bountyEl.classList.remove('hidden');
                    bountyEl.textContent = bounty;
                } else {
                    bountyEl.classList.add('hidden');
                    bountyEl.textContent = "";
                }
            }
        }

        const redCard = row.querySelector('.red-player');
        if (redCard) {
            const bountyEl = redCard.querySelector('.bounty-tag');
            if (bountyEl) {
                const bounty = redPlayers[rPlayer] || 0;
                if (bounty > 0) {
                    bountyEl.classList.remove('hidden');
                    bountyEl.textContent = bounty;
                } else {
                    bountyEl.classList.add('hidden');
                    bountyEl.textContent = "";
                }
            }
        }
    }
}

function updateStandingsTable(data) {
    if (!data) return;
    if (data.groupTitle) {
        const headerTitle = document.querySelector('.box-header h3');
        if (headerTitle) headerTitle.textContent = data.groupTitle.toUpperCase();
    }

    const tbody = document.querySelector('.standings-table tbody');
    if (!tbody) return;

    tbody.innerHTML = '';
    const teamsToRender = [];

    function tryAddTeam(teamData, specialClass = '') {
        if (teamData && teamData.tag) {
            teamsToRender.push({
                position: parseInt(teamData.position, 10) || 0,
                tag: teamData.tag,
                img: teamData.img || 'assets/images/Logos/square_default.jpg',
                wins: teamData.wins,
                losses: teamData.losses,
                points: teamData.points,
                cssClass: specialClass
            });
        }
    }

    tryAddTeam(data.leadingTeam);
    tryAddTeam(data.targetTeam, 'own-team');
    tryAddTeam(data.trailingTeam);
    
    teamsToRender.sort((a,b) => a.position - b.position);

    if (data.lastTeam && data.lastTeam.tag) {
        const lastTeamPos = parseInt(data.lastTeam.position, 10);
        const teamAbove = teamsToRender[teamsToRender.length - 1];
        const needsGap = teamAbove ? (lastTeamPos - teamAbove.position > 1) : false;

        tryAddTeam(data.lastTeam, 'enemy-team');
        
        if (needsGap) {
            teamsToRender[teamsToRender.length - 1].insertGapBefore = true;
        }
    }

    teamsToRender.sort((a,b) => a.position - b.position);

    teamsToRender.forEach(team => {
        if (team.insertGapBefore) {
            const gapRow = document.createElement('tr');
            gapRow.classList.add('table-gap');
            gapRow.innerHTML = `<td colspan="4"><div></div></td>`;
            tbody.appendChild(gapRow);
        }

        const tr = document.createElement('tr');
        if (team.cssClass) tr.classList.add(team.cssClass);

        tr.innerHTML = `
            <td class="pos-cell">${team.position}</td>
            <td class="team-cell">
                <img src="${team.img}" alt="${team.tag} Logo" class="team-logo">
                <span class="name-text">${team.tag}</span>
            </td>
            <td class="wl-cell">${team.wins} - ${team.losses}</td>
            <td class="pts-cell">${team.points}</td>
        `;
        tbody.appendChild(tr);
    });

    /*const gapRow = document.createElement('tr');
    gapRow.classList.add('table-gap');
    gapRow.innerHTML = `<td colspan="4"><div></div></td>`;
    tbody.appendChild(gapRow);
    */

    function formatEpoch(epochInSeconds) {
        const date = new Date(epochInSeconds * 1000);
        const formatter = new Intl.DateTimeFormat('de-DE', {
            day: 'numeric',
            month: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
            hour12: false,
            timeZone: 'Europe/Berlin'
        });
        return formatter.format(date).replace(',', '');
    }

    if (data.nextMatch) {
        const nextGameRow = document.createElement('tr');
        const date = formatEpoch(data.nextMatch.matchTime);
        nextGameRow.classList.add('next-game-row');
        nextGameRow.innerHTML = `
            <td colspan="4">
                <div class="next-game-container">
                    <span class="next-game-label">Next Game:</span>
                    <div class="next-game-team">
                        <img src="${data.nextMatch.img || 'assets/images/Logos/square_default.jpg'}" alt="Enemy Logo" class="team-logo">
                        <span class="name-text">${data.nextMatch.tag || 'TBD'}</span>
                    </div>
                    <span class="next-game-date">${date || 'open'}</span>
                </div>
            </td>
        `;
        tbody.appendChild(nextGameRow);
    } else {
        const nextGameRow = document.createElement('tr');
        nextGameRow.classList.add('next-game-row');
        nextGameRow.innerHTML = `
            <td colspan="4">
                <div class="next-game-container">
                    <span class="next-game-label">Next Game:</span>
                    <div class="next-game-team">
                        <img src="assets/images/Logos/square_default.jpg" alt="Enemy Logo" class="team-logo">
                        <span class="name-text">TBD</span>
                    </div>
                    <span class="next-game-date">open</span>
                </div>
            </td>
        `;
        tbody.appendChild(nextGameRow);
    }
}

function updateScoreTeams(data) {
    if (!data || !data.currentMatch) return;

    const scoreDisplay = document.getElementById('score-display');
    if (!scoreDisplay) return;

    const match = data.currentMatch;
    const op1 = match.opponent1;
    const op2 = match.opponent2;

    if (op1) {
        const blueSide = document.querySelector('.team-side.blue-side');
        if (blueSide) {
            const nameEl = blueSide.querySelector('.team-name');
            const recordEl = blueSide.querySelector('.team-record');
            if (nameEl) nameEl.textContent = op1.tag;
            if (recordEl) recordEl.textContent = `${op1.wins}-${op1.losses}`

            const logoEl = blueSide.querySelector('.score-team-logo');
            if (logoEl && op1.img) {
                logoEl.src = op1.img;
                logoEl.alt = `${op1.tag} Logo`;
            }

            const blueTicks = blueSide.querySelectorAll('.series-ticks .tick');
            blueTicks.forEach((tick, index) => {
                if (index < op1.matchScore) {
                    tick.classList.add('active');
                } else {
                    tick.classList.remove('active');
                }
            });
        }
    }

    if (op2) {
        const redSide = document.querySelector('.team-side.red-side');
        if (redSide) {
            const nameEl = redSide.querySelector('.team-name');
            const recordEl = redSide.querySelector('.team-record');
            if (nameEl) nameEl.textContent = op2.tag;
            if (recordEl) recordEl.textContent = `${op2.wins}-${op2.losses}`

            const logoEl = redSide.querySelector('.score-team-logo');
            if (logoEl && op2.img) {
                logoEl.src = op2.img;
                logoEl.alt = `${op2.tag} Logo`;
            }

            const redTicks = redSide.querySelectorAll('.series-ticks .tick');
            redTicks.forEach((tick, index) => {
                if (index < op2.matchScore) {
                    tick.classList.add('active');
                } else {
                    tick.classList.remove('active');
                }
            });
        }
    }
}

// ============================================================================
// INITIALIZATION
// ============================================================================
document.addEventListener('DOMContentLoaded', () => {
    // 1. WebSockets initialisieren
    connectSpectatorServer();
    connectPythonScript();

    // 2. Prime League HTTP API initialisieren & Interval starten
    fetchAndUpdateStandings();
    setInterval(fetchAndUpdateStandings, INTERVAL_ONE_HOUR);
});



//completed in div when completed

//                    <div class="summ-slot cd">
//                        <img src="assets/images/summs/cloud-drake.png" alt="F">
//                        <div class="cooldown-overlay"></div>

//                        <div class="ult-slot ready">
//                           <img src="assets/images/champs/cloud-drake.png" alt="U">
//                        </div>

//            <span class="spawn-time live">LIVE</span>

//                <div class="drakes-list blue-drakes">
//                    <img src="assets/images/ui/normal-drake/32px-Infernal_Dragon_Soul_buff - normal.png" class="drake-icon" alt="I">
//                    <img src="assets/images/ui/normal-drake/32px-Ocean_Dragon_Soul_buff - normal.png" class="drake-icon" alt="O">
//                </div>

//                    <span class="gold-lead blue-side-lead">+1.4K</span>

//            <div class="lane-gold-diff blue-lead">
//                <span class="diff-amount">0.0K</span>

/*                    <div class="ward-slot">
                        <img src="assets/images/items/cloud-drake.png" alt="W">
                        <span class="ward-count">3</span>
                    </div>
                    <div class="quest-slot"><img src="assets/images/items/cloud-drake.png" alt="Q"></div>
                </div>

                <div class="items-line">
                    <div class="item-slot">
                        <img src="assets/images/items/cloud-drake.png" alt="Black Cleaver">
                        <span class="item-count">2</span>
                    </div>
                    <div class="item-slot cd-item">
                        <img src="assets/images/items/cloud-drake.png" alt="BOTRK">
                        <div class="item-cooldown-overlay"></div>
                    </div>
                    <div class="item-slot empty"></div>
                    <div class="item-slot empty"></div>
                    <div class="item-slot empty"></div>
                    <div class="item-slot empty"></div>
                </div>
*/

/*
                <div class="summs-container">
                    <div class="summ-slot">
                        <img src="assets/images/summs/cloud-drake.png" alt="F">
                    </div>
                    <div class="summ-slot">
                        <img src="assets/images/summs/cloud-drake.png" alt="T">
                    </div>
                </div>
*/

/*
                <div class="runes-box">
                    <img src="assets/images/runes/cloud-drake.png" class="main-rune" alt="R">
                    <img src="assets/images/runes/cloud-drake.png" class="sub-rune" alt="R">
                </div>
*/

/*
                <div class="champ-avatar-box">
                    <img src="assets/images/champs/cloud-drake.png" class="champ-img" alt="Champ">
                    <div class="champ-level">5</div>
                    <div class="death-timer-overlay hidden"></div>
                    <div class="bounty-tag hidden"></div>
                </div>
*/

/*
                    <tbody>
                        <tr>
                            <td class="pos-cell">1</td>
                            <td class="team-cell">
                                <img src="assets/images/Logos/DRR.jpg" alt="" class="team-logo">
                                <span class="name-text">DRR</span>
                            </td>
                            <td class="wl-cell">3 - 0</td>
                            <td class="pts-cell">9</td>
                        </tr>
                        <tr class="own-team">
                            <td class="pos-cell">2</td>
                            <td class="team-cell">
                                <img src="assets/images/Logo-removebg-preview.png" alt="" class="team-logo">
                                <span class="name-text">AGG</span>
                            </td>
                            <td class="wl-cell">2 - 1</td>
                            <td class="pts-cell">6</td>
                        </tr>
                        <tr>
                            <td class="pos-cell">2</td>
                            <td class="team-cell">
                                <img src="assets/images/Logos/catgirls.jpg" alt="" class="team-logo">
                                <span class="name-text">LES</span>
                            </td>
                            <td class="wl-cell">2 - 1</td>
                            <td class="pts-cell">6</td>
                        </tr>
                        <tr class="table-gap">
                            <td colspan="4"><div></div></td>
                        </tr> 
                        <tr class="enemy-team">
                            <td class="pos-cell">8</td>
                            <td class="team-cell">
                                <img src="assets/images/Logos/square_default.jpg" alt="" class="team-logo">
                                <span class="name-text">SüFü</span>
                            </td>
                            <td class="wl-cell">0 - 3</td>
                            <td class="pts-cell">0</td>
                        </tr>
                        
                        <tr class="next-game-gap">
                            <td colspan="4"><div></div></td>
                        </tr>
                        
                        <tr class="next-game-row">
                            <td colspan="4">
                                <div class="next-game-container">
                                    <span class="next-game-label">Next Game:</span>
                                    <div class="next-game-team">
                                        <img src="assets/images/Logos/square_default.jpg" alt="Enemy Logo" class="team-logo">
                                        <span class="name-text">SüFü</span>
                                    </div>
                                    <span class="next-game-date">offen</span>
                                </div>
                            </td>
                        </tr>
                    </tbody>
*/