# DragonsNShit: World Crisis (Vertical Slice 0)

## Goal
Demonstrate a server-wide existential event that:

- is visible and legible to everyone
- cannot be solved by a single guild (even rich/stacked)
- creates multiple valid roles (combat + non-combat)
- has stateful progression, escalation, and outcome
- produces lasting world change (even if only one town/zone is affected in VS0)

**Non-goals (explicitly deferred)**

- multiple threats / rotating seasons
- complex economy simulation
- long quest chains / VO / cinematics beyond minimal
- full admin panel (config file + console prints are enough for VS0)

## Acceptance Criteria

### A) Threat Lifecycle

**A1. World-state machine exists**

- Event has clearly defined phases:
  - OMENS → BURROW → EMERGENCE → SPLIT WAR → FINAL WINDOW → RESOLUTION
- Server authoritative phase transitions based on:
  - timers, and
  - objective completion gates
- **Pass/Fail:** If event can be “forced” into kill state just by DPSing one entity, fail.

**A2. Global meter**

- A single global value exists (e.g., `LEY_INTEGRITY` 0..100) that:
  - decreases automatically over time during active phases
  - decreases faster when objectives fail
  - increases when server completes stabilizing objectives
- **Pass/Fail:** Meter must be visible in UI and queryable via server console/log.

**A3. Telegraphing is unavoidable**

At least 3 telegraphs must occur before full combat window:

1. skybox/audio/lighting change in the threatened region
2. world map indicator / “fault line” visualization
3. NPC/system message warning with countdown to escalation

**Pass/Fail:** A player logging in cold must notice within 60 seconds.

### B) “Rich Guild Can’t Solo” Hard Locks

**B1. Multi-objective simultaneity gate**

Final vulnerability window requires ≥3 objectives completed concurrently within a time window (e.g., 5 minutes):

- Anchor Objective (keep pylons alive)
- Ritual Objective (stabilize nodes)
- Intercept Objective (stop a roaming head / convoy / add wave)

**Pass/Fail:** If one group can sequentially do them alone, fail.

**B2. Contribution cap / diminishing returns**

Any single party/guild doing repeated identical actions gains reduced progress after threshold (server-side).

Diminishing returns applies to:

- damage to “progress-critical” targets
- anchor repair
- ritual contribution

**Pass/Fail:** A maxed guild stacking 40 people in one spot cannot brute-force the entire phase.

**B3. Geo-separation constraint**

- At least 2 critical objectives occur in separate locations requiring travel time or simultaneous teams.
- Travel cannot be bypassed by a single instant teleport mechanic in VS0.

**Pass/Fail:** If one zerg can rotate and still meet the concurrency gate, fail.

**B4. Anti-grief stability**

Griefing cannot permanently lock the event.

If players sabotage objectives, the system must:

- degrade gracefully
- still allow recovery with extra work
- log sabotage actions for admin review (VS0: simple log lines)

**Pass/Fail:** Event can be completed even with a hostile subgroup, given enough defenders.

### C) Roles and “Everyone Has a Job”

**C1. Minimum viable role diversity**

Event must include at least 4 distinct contribution roles:

- Strike (combat DPS/tank/heal)
- Builder (craft/transport/repair anchors)
- Ritualist (non-combat interaction/puzzle channel)
- Scout/Runner (intel + map tasks + delivery)

**Pass/Fail:** Each role must have a progression bar it can meaningfully advance.

**C2. Non-combat matters**

Builder + Ritualist contributions must be strictly required to reach Final Window.

**Pass/Fail:** If combat-only groups can complete all objectives, fail.

### D) Combat Design Requirements

**D1. Boss is not a “single health bar”**

The threat has:

- an invulnerable state
- at least one “break armor / expose weak point” mechanic
- add waves that threaten objectives (anchors/ritualists)

**Pass/Fail:** If it’s just “hit it until it dies,” fail.

**D2. Split War exists in VS0**

- At least two sub-bosses/heads with different mechanics spawn in different locations.
- Ignoring one makes final window harder (e.g., shorter vulnerability window, higher add pressure).

**Pass/Fail:** Players must make a strategic choice and feel consequences.

### E) Persistence and Consequences

**E1. Outcome changes the world**

On victory:

- unlock a visible world change (e.g., ley bridge, safe corridor, monument, vendor unlocked)

On failure:

- apply a scar to one region/town (e.g., debuff zone, closed services, hostile spawns)

**Pass/Fail:** State must persist across server restart.

**E2. Event cooldown**

- Event cannot be immediately retriggered.
- Cooldown is server-configurable (VS0: config constant).

**Pass/Fail:** If admins restart the server and it retriggers instantly, fail unless config says so.

### F) Scoring, Rewards, and Fairness

**F1. Multi-axis merit**

Merit earned from at least 5 sources:

- boss damage
- healing/defending anchor
- crafting/delivery
- ritual progress actions
- scouting/discovery/delivery completions

**Pass/Fail:** UI shows personal merit + category breakdown.

**F2. Rewards are tiered, not rank-1**

- Rewards based on merit brackets (e.g., Bronze/Silver/Gold) rather than “top 10 only”.
- At least one cosmetic reward + one functional reward exist in VS0.

**Pass/Fail:** A casual Builder can earn a meaningful reward without touching DPS.

### G) Observability + Admin Controls

**G1. Server observability**

Server outputs periodic status line:

- current phase
- time to next escalation
- global meter
- objective states (anchors alive count, ritual progress %, head status)

**Pass/Fail:** Debugging does not require attaching a debugger.

**G2. Config knobs (VS0 minimum)**

- `event_enabled`
- `phase_durations`
- `meter_decay_rate`
- `required_concurrent_objectives`
- `pvp_during_event` (on/off)

**Pass/Fail:** Changing knobs changes behavior without code edits.

## Definition of Done

### 1) Playtest Proofs (must be reproducible)

**DoD-1: Solo + small-group failure proof**

Run with 1 party (4–6 players) + rich gear:

- they can participate
- they cannot reach Final Window alone

Evidence: recorded log + in-game outcome.

**DoD-2: Rich guild brute-force proof**

Run with 1 large coordinated guild (e.g., 20–40):

- if they stack one objective, diminishing returns kicks in
- geo-separation forces split teams
- concurrency gate blocks sequential completion

Evidence: server logs show gate enforcement.

**DoD-3: Mixed-server success proof**

Run with mixed random players + one guild:

- the event can be completed
- at least 3 different roles contribute materially

Evidence: merit breakdown shows non-combat roles in top brackets.

### 2) Technical Completion

**DoD-4: Server-authoritative state**

- All phase transitions + objective completion are computed server-side.
- Clients cannot spoof progress via packets. (VS0: sanity checks + server-owned timers)

**DoD-5: Persistence**

Event outcome persists across restart:

- phase resets safely
- scar/victory changes remain
- cooldown remains

**DoD-6: Performance sanity**

- Event running with max expected players does not degrade below target tickrate beyond agreed threshold.
- Objective interactions do not allocate unbounded memory.

### 3) Content Completion (VS0 minimum)

**DoD-7: One complete threat**

One threat fully playable end-to-end with:

- at least 2 regions involved
- at least 2 sub-bosses/heads
- 3 objective types
- win + fail outcome

**DoD-8: UX**

Players can answer within 30 seconds:

- “What phase are we in?”
- “What do I do to help?”
- “Where do I go?”

Evidence: UI elements present + tested by a fresh player.

### 4) Anti-cheese / Anti-degenerate strategies checklist

**DoD-9: Cheese list validated**

The following must be tested and blocked or absorbed:

- stacking 40 players on one objective
- kiting boss out of bounds
- ignoring all heads and still winning
- repeatedly farming merit via one low-risk action
- griefers trying to hard-lock progress

---

## World Crisis VS0 Implementation Contract

This is the server-facing contract for VS0. If these are not true, VS0 fails.

### 0) Hard invariants

- Server-authoritative crisis state: phases, meters, objectives, scoring computed server-side only.
- Concurrency gate: final vulnerability requires N objectives in distinct regions completed within the same window.
- Geo separation: objectives are physically separated and cannot be “one zerg rotates.”
- Diminishing returns: repeated contribution from the same identity group saturates progress.
- Multi-role requirement: at least one objective is non-combat and required.

### 1) Data model (server)

Add this to a server-only header (preferred) or `protocol.h` if unavoidable.

```c
typedef enum {
  CRISIS_OFF=0,
  CRISIS_OMENS,
  CRISIS_BURROW,
  CRISIS_EMERGENCE,
  CRISIS_SPLIT_WAR,
  CRISIS_FINAL_WINDOW,
  CRISIS_RESOLUTION
} CrisisPhase;

typedef enum {
  OBJ_ANCHOR=0,
  OBJ_RITUAL,
  OBJ_INTERCEPT,
  OBJ_COUNT
} CrisisObjType;

typedef enum {
  OBJ_INACTIVE=0,
  OBJ_ACTIVE,
  OBJ_COMPLETE,
  OBJ_FAILED
} CrisisObjState;

typedef struct {
  int id;
  CrisisObjType type;
  CrisisObjState state;

  // location/region identity
  int region_id;
  float x, y, z;

  // progress
  float progress;        // 0..1
  float progress_rate;   // tuning
  unsigned int last_tick;

  // gating & timing
  unsigned int activated_ms;
  unsigned int completed_ms;
  unsigned int fail_deadline_ms;

  // anti-solo: contributions by cohort
  // small fixed ring map: last K cohorts that contributed.
  // cohort_id = hash(guild_id or party_id or player_id fallback)
  unsigned int cohort_id[16];
  float cohort_contrib[16];    // normalized 0..1
} CrisisObjective;

typedef struct {
  int active;
  int seed;
  CrisisPhase phase;

  // global meter
  float ley_integrity; // 0..100
  float decay_per_sec;

  // phase timers
  unsigned int phase_started_ms;
  unsigned int next_escalation_ms;

  // objectives
  CrisisObjective obj[OBJ_COUNT];

  // concurrency gate
  int required_concurrent;           // e.g. 3
  unsigned int concurrency_window_ms; // e.g. 300000 (5 min)
  unsigned int last_gate_open_ms;
  int gate_open;                     // whether final window is open

  // heads/subbosses minimal
  int head_alive[2];                 // VS0: 2 heads
  int head_region[2];

  // merit
  // VS0: just per-client totals + per-category totals
  float merit_total[MAX_CLIENTS];
  float merit_cat[MAX_CLIENTS][8]; // dmg/heal/build/ritual/scout/defense/etc

  // outcome persistence
  int last_outcome;                  // 0 none, 1 win, 2 fail
  unsigned int cooldown_until_ms;
  int scar_region_id;
} CrisisState;
```

Where it lives: add `CrisisState crisis;` to `ServerState` server-side only, or a new `server_crisis.h`.

### 2) Identity / “Rich guild can’t solo” anti-bruteforce

You need a stable cohort key. VS0 rule:

- If you have guilds: `cohort_id = guild_id`
- Else parties: `cohort_id = party_id`
- Else fallback: `cohort_id = client_id`

Diminishing returns function:

```c
float cohort_scale(float c) {
  // c = cohort_contrib for this objective, 0..inf
  // fast saturation so one cohort can’t do 100% alone
  // example: first 20% at full, then drop hard
  if (c < 0.2f) return 1.0f;
  if (c < 0.5f) return 0.35f;
  return 0.10f;
}
```

This is the “rich guild can still help a ton, but cannot complete everything alone” lever.

### 3) Objective types (VS0 minimal mechanics)

**OBJ_ANCHOR (Builder + defenders)**  
World has an “Anchor Pylon” object with HP and repair interaction.

Progress is “anchors stabilized” (0..1) based on:

- pylon alive + repaired above threshold
- time held under pressure

Mechanic:

- Every tick, if pylon HP > 70%: progress += rate * dt
- Add waves attack it (simple spawn near objective)
- Builders can repair with delivered kits (or hold-interact in VS0)

Non-combat requirement: if repair kits exist, only Builders can meaningfully keep it alive.

**OBJ_RITUAL (Ritualists + defenders)**  
A channel interaction with “runes” (VS0: 3 nodes).

Players must hold E on 3 runes simultaneously for 10 seconds total.

Implementation:

- Track rune_active_count each tick
- If >= 3: progress += rate * dt else decay slightly
- Damage/knockback interrupts channel

This is pure “coordination over DPS”.

**OBJ_INTERCEPT (Strike + scouts)**  
A roaming “Convoy/Head Spawn” target in a different region.

Progress increases when:

- “intercept boss” takes armor segments down
- or convoy reaches 0 HP before a timer

This forces distribution: you can’t keep everyone on pylons.

### 4) Crisis phase flow (server tick state machine)

Runs inside server main loop at ~60Hz.

**Phase transitions**

**OMENS**

- Telemetry + UI only, meter stable or slow decay
- After `T_omens` → `BURROW`

**BURROW**

- start meter decay
- spawn objective markers, but objectives not active yet
- after `T_burrow` → `EMERGENCE`

**EMERGENCE**

- activate objectives: Anchor + Ritual + Intercept in distinct regions
- require players to start progress
- if any objective fails hard: meter decay increases
- when enough objectives complete at least once → `SPLIT_WAR`

**SPLIT_WAR**

- activate 2 heads in different regions (VS0)
- each head modifies final window if ignored:
  - head A alive reduces final window duration
  - head B alive increases add pressure
- when heads are “weakened” (HP below threshold) and objectives are completable → allow gate checks

**FINAL_WINDOW**

- opens only if concurrency condition is met
- worm/boss becomes damageable for limited time
- if boss dies → `RESOLUTION` (win)
- if timer runs out → fallback to `SPLIT_WAR` (harder) or fail depending on `ley_integrity`

**RESOLUTION**

- apply outcome world change + cooldown
- reset active combat spawns, keep scar/monument state

**Concurrency gate check**

Every tick (or once per second), compute:

- count objectives that are `OBJ_COMPLETE` with `completed_ms >= now - concurrency_window_ms`
- ensure they are in distinct regions (or at least 2 distinct regions in VS0)

If count >= `required_concurrent`:

- `gate_open = 1`
- `last_gate_open_ms = now`
- transition to `CRISIS_FINAL_WINDOW`

### 5) Boss vulnerability model (VS0)

You don’t need a real boss AI yet. You need invulnerability + window.

- Boss has `invuln = 1` until `FINAL_WINDOW`.
- In `FINAL_WINDOW`:
  - `invuln = 0`
  - boss HP drains from combat
  - anchor/ritual/intercept add pressure continues
  - fail condition: if boss HP not reduced to 0 before `final_window_ms`, revert to `SPLIT_WAR` and reduce `ley_integrity`.

### 6) Merit scoring (multi-axis)

VS0: keep it dead simple, but multi-axis.

**Categories (suggested indices)**

0 dmg_to_boss  
1 dmg_to_objective_threats (adds/heads)  
2 healing_or_support (if exists; else “defense time near objective”)  
3 building/repairs/deliveries  
4 ritual_time  
5 scouting/interaction completions  
6 pvp_defense (optional)  
7 misc

**Award rules**

- Damage merit: `+damage * k1`
- Ritual: `+dt_held * k2`
- Build: `+repair_amount * k3` or `+delivery_count * k4`
- Intercept: `+armor_breaks * k5`
- Defense: `+seconds_in_objective_radius_while_objective_active * k6`

**Bracket rewards (server computed)**

- Bronze: `merit_total >= X`
- Silver: `merit_total >= Y`
- Gold: `merit_total >= Z`

### 7) Persistence (VS0 minimum)

Write a tiny file `crisis_state.dat` (binary or JSON-ish) that stores:

- `cooldown_until_ms` (or remaining seconds)
- `last_outcome`
- `scar_region_id` and any “scar flags”
- optional seed so event reproducible

Load at server boot, apply immediately.

### 8) Client UX (VS0 minimum)

Client must show within 30 seconds:

- current phase name
- global meter (e.g., `LEY INTEGRITY: 72%`)
- objective list with location labels:
  - `ANCHOR: BONEYARD EAST (42%)`
  - `RITUAL: SKYRIDGE (ACTIVE)`
  - `INTERCEPT: DUST ROAD (COMPLETE)`
- If `FINAL_WINDOW` open: big banner `VULNERABLE: 02:34`

No fancy UI needed; draw text in HUD.

### 9) Test checklist (the “prove Avengers” list)

- One guild zerg: stack everyone on Anchor → progress slows after 20% and can’t finish without Ritual + Intercept simultaneously.
- Sequential cheese: complete Anchor, then run to Ritual, then Intercept → gate never opens (concurrency window expires).
- Non-combat necessity:
  - remove all fighters, keep builders/ritualists → objectives progress but can’t defend (fails).
  - remove builders/ritualists, keep fighters → cannot reach FINAL_WINDOW (fails).
- Recovery from sabotage: griefers kill pylons → defenders can rebuild and continue (not hard-locked).
- Persistence: restart server mid-event → comes up in safe state (either resumes with timers clamped or resets phase conservatively), scar/outcome persists.
