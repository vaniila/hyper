package logger

import "log"

type logger struct {
	id string
}

// Start logger server (unused by default)
func (l *logger) Start() error {
	return nil
}

// Stop logger server (unused by default)
func (l *logger) Stop() error {
	return nil
}

// Debug log at debug level
func (l *logger) Debug(msg string, fields ...Field) {
	if len(fields) == 0 {
		log.Print(msg)
		return
	}
	log.Printf(msg, l.fields(fields...))
}

// Info log at info level
func (l *logger) Info(msg string, fields ...Field) {
	if len(fields) == 0 {
		log.Print(msg)
		return
	}
	log.Printf(msg, l.fields(fields...))
}

// Warn log at warn level
func (l *logger) Warn(msg string, fields ...Field) {
	if len(fields) == 0 {
		log.Print(msg)
		return
	}
	log.Printf(msg, l.fields(fields...))
}

// Error log at error level
func (l *logger) Error(msg string, fields ...Field) {
	if len(fields) == 0 {
		log.Print(msg)
		return
	}
	log.Printf(msg, l.fields(fields...))
}

// Fatal log at fatal level
func (l *logger) Fatal(msg string, fields ...Field) {
	if len(fields) == 0 {
		log.Fatal(msg)
		return
	}
	log.Fatal(msg, l.fields(fields...))
}

// Panic log at panic level
func (l *logger) Panic(msg string, fields ...Field) {
	if len(fields) == 0 {
		log.Panic(msg)
		return
	}
	log.Panic(msg, l.fields(fields...))
}

// String returns the logger name
func (l *logger) String() string {
	return "Hyper::Logger"
}

// fields convert fields to a slice of interfaces (used for default logger)
func (l *logger) fields(fields ...Field) []interface{} {
	var is = make([]interface{}, len(fields))

	for _, f := range fields {
		is = append(is, f.Value())
	}

	return is
}
