# Raksana Backend - Structured Entity Relationship Diagram

This document contains the database schema organized by functional domains to provide better understanding of the system architecture.

## 📊 Database Overview

The Raksana backend database is organized into the following domains:
- **User Domain**: User management, profiles, statistics, and logs
- **Habit Domain**: Habit tracking, packets, tasks, and recaps
- **Activity Domain**: Challenges, quests, events, and treasures
- **Sustainability Domain**: Environmental scanning, greenprints, and eco-friendly activities
- **System Domain**: Sessions and system-level tables

---

## 👤 User Domain

### Core User Management
```sql
Table "users" {
  "id" bigint [pk, not null]
  "name" "character varying(255)" [not null]
  "username" "character varying(255)" [unique, not null]
  "email" "character varying(255)" [unique, not null]
  "email_verified_at" timestamp(0)
  "password" "character varying(255)" [not null]
  "is_admin" boolean [not null, default: false]
  "remember_token" "character varying(100)"
  "created_at" timestamp(0)
  "updated_at" timestamp(0)
}

Table "profiles" {
  "id" bigint [pk, not null]
  "user_id" bigint [not null]
  "current_exp" bigint [not null, default: `'0'::bigint`]
  "exp_needed" bigint [not null, default: `'100'::bigint`]
  "level" integer [not null, default: 1]
  "points" bigint [not null, default: `'0'::bigint`]
  "profile_key" "character varying(255)" [not null, default: `'profiles/Portrait_Placeholder.png'::charactervarying`]
}
```

### User Statistics & Analytics
```sql
Table "statistics" {
  "id" bigint [pk, not null]
  "user_id" bigint [not null]
  "challenges" integer [not null, default: 0]
  "events" integer [not null, default: 0]
  "quests" integer [not null, default: 0]
  "treasures" integer [not null, default: 0]
  "longest_streak" integer [not null, default: 0]
  "tree_grown" integer [not null, default: 0]
}

Table "histories" {
  "id" bigint [pk, not null]
  "user_id" bigint [not null]
  "name" "character varying(255)" [not null]
  "type" "character varying(255)" [not null]
  "category" "character varying(255)" [not null]
  "amount" integer [not null]
  "created_at" timestamp(0) [not null, default: `CURRENT_TIMESTAMP`]
}
```

### User Logs & Memories
```sql
Table "logs" {
  "id" bigint [pk, not null]
  "user_id" bigint [not null]
  "text" text [not null]
  "is_system" boolean [not null, default: false]
  "is_private" boolean [not null]
  "created_at" timestamp(0) [not null, default: `CURRENT_TIMESTAMP`]
}

Table "memories" {
  "id" bigint [pk, not null]
  "user_id" bigint [not null]
  "file_key" "character varying(255)" [not null]
  "description" text [not null]
  "created_at" timestamp(0) [not null, default: `CURRENT_TIMESTAMP`]
}
```

---

## 🎯 Habit Domain

### Habit Management
```sql
Table "packets" {
  "id" bigint [pk, not null]
  "user_id" bigint [not null]
  "name" "character varying(255)" [not null]
  "target" "character varying(255)" [not null]
  "description" text [not null]
  "completed_task" integer [not null, default: 0]
  "expected_task" integer [not null]
  "task_per_day" integer [not null]
  "completed" boolean [not null, default: false]
  "created_at" timestamp(0) [not null, default: `CURRENT_TIMESTAMP`]
}

Table "habits" {
  "id" bigint [pk, not null]
  "packet_id" bigint [not null]
  "name" "character varying(255)" [not null]
  "description" text [not null]
  "difficulty" "character varying(255)" [not null]
  "locked" boolean [not null]
  "weight" integer [not null]
}

Table "tasks" {
  "id" bigint [pk, not null]
  "habit_id" bigint [not null]
  "user_id" bigint [not null]
  "packet_id" bigint [not null]
  "name" "character varying(255)" [not null]
  "description" text [not null]
  "difficulty" "character varying(255)" [not null]
  "completed" boolean [not null, default: false]
  "created_at" timestamp(0) [not null, default: `CURRENT_TIMESTAMP`]
  "updated_at" timestamp(0)
}
```

### Habit Recaps & Analytics
```sql
Table "recaps" {
  "id" bigint [pk, not null]
  "user_id" bigint [not null]
  "summary" text [not null]
  "tips" text [not null]
  "assigned_task" integer [not null]
  "completed_task" integer [not null]
  "completion_rate" "character varying(255)" [not null]
  "growth_rating" "character varying(255)" [not null]
  "type" "character varying(255)" [not null]
  "created_at" timestamp(0) [not null, default: `CURRENT_TIMESTAMP`]
}

Table "recap_details" {
  "id" bigint [pk, not null]
  "monthly_recap_id" bigint [not null]
  "challenges" integer [not null, default: 0]
  "events" integer [not null, default: 0]
  "quests" integer [not null, default: 0]
  "treasures" integer [not null, default: 0]
  "longest_streak" integer [not null, default: 0]
  "created_at" timestamp(0) [not null, default: `CURRENT_TIMESTAMP`]
}
```

---

## 🏃‍♂️ Activity Domain

### Activity Details & Base
```sql
Table "details" {
  "id" bigint [pk, not null]
  "name" "character varying(255)" [not null]
  "description" text [not null]
  "point_gain" bigint [not null]
  "created_at" timestamp(0)
  "updated_at" timestamp(0)
}

Table "codes" {
  "id" "character varying(255)" [pk, not null]
  "image_url" "character varying(255)" [not null]
}
```

### Challenges
```sql
Table "challenges" {
  "id" bigint [pk, not null]
  "detail_id" bigint [not null]
  "day" integer [not null]
  "difficulty" "character varying(255)" [not null]
}

Table "participations" {
  "id" bigint [pk, not null]
  "challenge_id" bigint [not null]
  "user_id" bigint [not null]
  "memory_id" bigint [not null]
  "created_at" timestamp(0) [not null, default: `CURRENT_TIMESTAMP`]
}
```

### Quests
```sql
Table "quests" {
  "id" bigint [pk, not null]
  "detail_id" bigint [not null]
  "code_id" "character varying(255)" [not null]
  "location" text [not null]
  "latitude" numeric(10,7)
  "longitude" numeric(10,7)
  "max_contributors" integer [not null]
  "finished" boolean [not null, default: false]
  "clue" text
}

Table "contributions" {
  "id" bigint [pk, not null]
  "quest_id" bigint [not null]
  "user_id" bigint [not null]
  "created_at" timestamp(0) [not null, default: `CURRENT_TIMESTAMP`]
}
```

### Events
```sql
Table "events" {
  "id" bigint [pk, not null]
  "detail_id" bigint [not null]
  "code_id" "character varying(255)" [not null]
  "location" text [not null]
  "latitude" numeric(10,7)
  "longitude" numeric(10,7)
  "contact" "character varying(255)" [not null]
  "starts_at" timestamp(0) [not null]
  "ends_at" timestamp(0) [not null]
  "cover_key" "character varying(255)"
}

Table "attendances" {
  "id" bigint [pk, not null]
  "user_id" bigint [not null]
  "event_id" bigint [not null]
  "attended" boolean [not null, default: false]
  "contact_number" "character varying(255)" [not null]
  "created_at" timestamp(0) [not null, default: `CURRENT_TIMESTAMP`]
  "attended_at" timestamp(0)
}
```

### Treasures
```sql
Table "treasures" {
  "id" bigint [pk, not null]
  "name" "character varying(255)" [not null]
  "point_gain" bigint [not null]
  "code_id" "character varying(255)" [not null]
  "claimed" boolean [not null, default: false]
  "created_at" timestamp(0)
  "updated_at" timestamp(0)
}

Table "claimed" {
  "id" bigint [pk, not null]
  "user_id" bigint [not null]
  "treasure_id" bigint [not null]
  "created_at" timestamp(0) [not null, default: `CURRENT_TIMESTAMP`]
}
```

---

## 🌱 Sustainability Domain

### Environmental Scanning
```sql
Table "scans" {
  "id" bigint [pk, not null]
  "user_id" bigint [not null]
  "title" "character varying(255)" [not null]
  "description" text [not null]
  "image_key" "character varying(255)" [not null]
  "created_at" timestamp(0) [not null, default: `CURRENT_TIMESTAMP`]
}

Table "items" {
  "id" bigint [pk, not null]
  "user_id" bigint [not null]
  "scan_id" bigint [not null]
  "name" "character varying(255)" [not null]
  "description" text [not null]
  "value" "character varying(255)" [not null]
  "created_at" timestamp(0) [not null, default: `CURRENT_TIMESTAMP`]
}
```

### Greenprints & Sustainability Guides
```sql
Table "greenprints" {
  "id" bigint [pk, not null]
  "item_id" bigint [not null]
  "image_key" "character varying(255)" [not null]
  "title" "character varying(255)" [not null]
  "description" text [not null]
  "sustainability_score" "character varying(255)" [not null]
  "estimated_time" "character varying(255)" [not null]
  "created_at" timestamp(0) [not null, default: `CURRENT_TIMESTAMP`]
}

Table "steps" {
  "id" bigint [pk, not null]
  "greenprint_id" bigint [not null]
  "description" text [not null]
  "created_at" timestamp(0) [not null, default: `CURRENT_TIMESTAMP`]
}

Table "materials" {
  "id" bigint [pk, not null]
  "name" "character varying(255)" [not null]
  "description" text [not null]
  "price" integer [not null]
  "quantity" integer [not null]
  "greenprint_id" bigint [not null]
}

Table "tools" {
  "id" bigint [pk, not null]
  "greenprint_id" bigint [not null]
  "name" "character varying(255)" [not null]
  "description" text [not null]
  "price" integer [not null]
  "created_at" timestamp(0) [not null, default: `CURRENT_TIMESTAMP`]
}
```

### Environmental Regions
```sql
Table "regions" {
  "id" bigint [pk, not null]
  "name" "character varying(255)" [not null]
  "location" text [not null]
  "latitude" doubleprecision [not null]
  "longitude" doubleprecision [not null]
  "tree_amount" integer [not null, default: 0]
  "created_at" timestamp(0)
  "updated_at" timestamp(0)
}
```

---

## 🔧 System Domain

### Session Management
```sql
Table "sessions" {
  "id" "character varying(255)" [pk, not null]
  "user_id" bigint
  "ip_address" "character varying(45)"
  "user_agent" text
  "payload" text [not null]
  "last_activity" integer [not null]

  Indexes {
    last_activity [type: btree, name: "sessions_last_activity_index"]
    user_id [type: btree, name: "sessions_user_id_index"]
  }
}
```

---

## 🔗 Relationships

### User Domain Relationships
```sql
// User Profile & Statistics
Ref "profiles_user_id_foreign":"users"."id" < "profiles"."user_id"
Ref "statistics_user_id_foreign":"users"."id" < "statistics"."user_id"
Ref "histories_user_id_foreign":"users"."id" < "histories"."user_id"
Ref "logs_user_id_foreign":"users"."id" < "logs"."user_id"
Ref "memories_user_id_foreign":"users"."id" < "memories"."user_id"
```

### Habit Domain Relationships
```sql
// Habit Management Chain
Ref "packets_user_id_foreign":"users"."id" < "packets"."user_id"
Ref "habits_packet_id_foreign":"packets"."id" < "habits"."packet_id"
Ref "tasks_habit_id_foreign":"habits"."id" < "tasks"."habit_id"
Ref "tasks_packet_id_foreign":"packets"."id" < "tasks"."packet_id"
Ref "tasks_user_id_foreign":"users"."id" < "tasks"."user_id"

// Recap System
Ref "recaps_user_id_foreign":"users"."id" < "recaps"."user_id"
Ref "recap_details_monthly_recap_id_foreign":"recaps"."id" < "recap_details"."monthly_recap_id"
```

### Activity Domain Relationships
```sql
// Challenge System
Ref "challenges_detail_id_foreign":"details"."id" < "challenges"."detail_id"
Ref "participations_challenge_id_foreign":"challenges"."id" < "participations"."challenge_id"
Ref "participations_memory_id_foreign":"memories"."id" < "participations"."memory_id"
Ref "participations_user_id_foreign":"users"."id" < "participations"."user_id"

// Quest System
Ref "quests_detail_id_foreign":"details"."id" < "quests"."detail_id"
Ref "quests_code_id_foreign":"codes"."id" < "quests"."code_id"
Ref "contributions_quest_id_foreign":"quests"."id" < "contributions"."quest_id"
Ref "contributions_user_id_foreign":"users"."id" < "contributions"."user_id"

// Event System
Ref "events_detail_id_foreign":"details"."id" < "events"."detail_id"
Ref "events_code_id_foreign":"codes"."id" < "events"."code_id"
Ref "attendances_event_id_foreign":"events"."id" < "attendances"."event_id"
Ref "attendances_user_id_foreign":"users"."id" < "attendances"."user_id"

// Treasure System
Ref "treasures_code_id_foreign":"codes"."id" < "treasures"."code_id"
Ref "claimed_treasure_id_foreign":"treasures"."id" < "claimed"."treasure_id"
Ref "claimed_user_id_foreign":"users"."id" < "claimed"."user_id"
```

### Sustainability Domain Relationships
```sql
// Scanning & Items
Ref "scans_user_id_foreign":"users"."id" < "scans"."user_id"
Ref "items_scan_id_foreign":"scans"."id" < "items"."scan_id"
Ref "items_user_id_foreign":"users"."id" < "items"."user_id"

// Greenprint System
Ref "greenprints_item_id_foreign":"items"."id" < "greenprints"."item_id"
Ref "steps_greenprint_id_foreign":"greenprints"."id" < "steps"."greenprint_id"
Ref "materials_greenprint_id_foreign":"greenprints"."id" < "materials"."greenprint_id"
Ref "tools_greenprint_id_foreign":"greenprints"."id" < "tools"."greenprint_id"
```

---

## 📈 Domain Interaction Summary

### Cross-Domain Relationships:
1. **User ↔ Habit**: Users own packets, habits, and tasks
2. **User ↔ Activity**: Users participate in challenges, contribute to quests, attend events, and claim treasures
3. **User ↔ Sustainability**: Users perform scans and create items that generate greenprints
4. **Activity ↔ User Memory**: Challenge participations reference user memories
5. **Habit ↔ User Statistics**: Habit completion affects user statistics and recaps

### Key Integration Points:
- **users** table serves as the central hub connecting all domains
- **details** table provides base information for activities (challenges, quests, events)
- **codes** table links activities (quests, events, treasures) with QR/code functionality
- **memories** table bridges user content with activity participation
