interface Lesson {
	UUID: string;
	TimetableEntryUUID: string;
	Date: Date;
	Cancelled: boolean;
}

interface LessonCooked {
	UUID: string;
	TimetableEntryUUID: string;
	Date: Date;
	Cancelled: boolean;
	Name: string;
	Position: number;
	StartTime: Date;
	EndTime: Date;
}

interface TimetableEntry {
	UUID: string;
	timetableUUID: string;
	subjectUUID: string;
	teacherUUID?: string;
	periodUUID: string;

	place: string;
	DayOfWeek: number;
}

interface Period {
	UUID: string;
	Name: string;
	Position: number;
	StartTime: Date;
	EndTime: Date;
}

interface Subject {
	UUID: string;
	Name: string;
}

export type { Lesson, LessonCooked, TimetableEntry, Period, Subject };

//	UUID               uuid.UUID `bun:"uuid,pk,type:uuid"`
//	TimetableEntryUUID uuid.UUID `bun:"timetable_entry_uuid,type:uuid,unique"`
//	Date               time.Time `bun:"date,notnull,unique"`
//	Cancelled          bool      `bun:"cancelled,type:boolean,default:false"`
/*     string UUID = 1;
    string timetableUUID = 2;
    string periodUUID = 3;
    string subjectUUID = 4;
    optional string teacherUUID = 5;
    
    string place = 6;
    DayOfWeek day_of_week = 7; */
/* string UUID = 1;
    string name = 2;
    int32 position = 3;
    google.protobuf.Timestamp start_time = 4;
    google.protobuf.Timestamp end_time = 5; */
