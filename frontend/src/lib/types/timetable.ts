enum Day {
	Monday = 0,
	Tuesday,
	Wednesday,
	Thursday,
	Friday,
	Saturday,
	Sunday
}

interface AbstractLesson {
	uuid: string;
	name: string;
}

interface ConcreteLesson {
	start: string;
	end: string;
	lessonInfo: AbstractLesson;
	teacherUuid: string;
	homeworkTextGetUrl?: string;
	homeworkAttachmentsGetUrls: string[];
	lessonUuid: string;
	marks: string[];
}

interface ConcreteLessonEntry {
	day: Day;
	lessons: ConcreteLesson[];
}

type Timetable = ConcreteLessonEntry[];

export type { Day, AbstractLesson, ConcreteLesson, ConcreteLessonEntry, Timetable };
