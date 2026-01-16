CREATE TABLE IF NOT EXISTS psychologists (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    phone TEXT UNIQUE,
    cancellation_window_hours INTEGER DEFAULT 24,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS patients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    psychologist_id INTEGER NOT NULL,
    email TEXT NOT NULL UNIQUE,
    phone TEXT UNIQUE,
    active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (psychologist_id) REFERENCES psychologists(id)
);

CREATE TABLE IF NOT EXISTS schedule_configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    psychologist_id INTEGER NOT NULL,
    day_of_week INTEGER NOT NULL,
    start_time TEXT NOT NULL,
    end_time TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (psychologist_id) REFERENCES psychologists(id),
    UNIQUE(psychologist_id, day_of_week, start_time)
);

CREATE TABLE IF NOT EXISTS recurring_slots (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    psychologist_id INTEGER NOT NULL,
    patient_id INTEGER NOT NULL,
    day_of_week INTEGER NOT NULL,
    start_time TEXT NOT NULL,
    duration_minutes INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (psychologist_id) REFERENCES psychologists(id),
    FOREIGN KEY (patient_id) REFERENCES patients(id),
    UNIQUE(psychologist_id, day_of_week, start_time)
);

CREATE TABLE IF NOT EXISTS appointments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    psychologist_id INTEGER NOT NULL,
    patient_id INTEGER NOT NULL,
    date DATE NOT NULL,
    start_time TEXT NOT NULL,
    duration_minutes INTEGER NOT NULL,
    status TEXT CHECK(status IN ('scheduled', 'cancelled', 'completed')) DEFAULT 'scheduled',
    rescheduled_from_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (psychologist_id) REFERENCES psychologists(id),
    FOREIGN KEY (patient_id) REFERENCES patients(id),
    FOREIGN KEY (rescheduled_from_id) REFERENCES appointments(id),
    UNIQUE(psychologist_id, date, start_time)
);

-- ÍNDICES DE RENDIMIENTO
CREATE INDEX IF NOT EXISTS idx_appointments_calendar 
ON appointments(psychologist_id, date);

CREATE INDEX IF NOT EXISTS idx_patients_psychologist 
ON patients(psychologist_id);