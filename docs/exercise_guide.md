# Exercise

For each lesson, there is an exercise that has several questions. A user can do the exercise multiple times, but the exercise will be regenerated every time the user wants to do again (versioning). The system will save the maximum score from all atempts of each user.

## Object

```typescript
class Exercise {
    _id: ObjectId;
    lessonId: ObjectId;
    version: number;
    questions: []Question;
    score: number;
    createdAt: Date;
    submittedAt: Date;
}

class Question {
    _id: ObjectId;
    exerciseId: ObjectId;
    question: string;
    choices: []string;
    answer: string;
    correctAnswer: string;
    score: number;
}
```

## API

### `POST /exercises/:lessonId`

- body
```json
{
    "questions": [
        {
            "question": string,
            "choices": []string,
            "correctAnswer": string,
            "score": number
        }
    ]
}
```

- Create new exercise by lessonId, version will be increased.
- For Frontend: Call this endpoint when clicking "ทำแบบทดสอบ" or "ทำใหม่".

### `POST /exercises/submit/:exerciseId`

- body
```json
{
    "answers": [
        {
            "questionId": ObjectId,
            "answer": string
        }
    ]
}
```

- Create attempt and calculate score by exerciseId.
- For Frontend: Call this endpoint when clicking "ส่งคำตอบ".

### `GET /exercises/:exerciseId`

- Get exercise by exerciseId
- For Frontend: Call this at the exercise page.

### `GET /exercises/lesson/:lessonId`

- Get all exercises by lessonId.
- For Frontend: Call this at the attempts history page. (We haven't had this page in the Figma yet)

### `GET /exercises/max-score/:lessonId`

- Get the maximum score from all atempts.
- For Frontend: Call this at the dashboard page for each lesson.

### `GET /exercises/average-score/:courseId`

- Get average score from the maximum score of each lesson by courseId
- For Frontend: Call this at the mentor dashboard page when mentors look at the dashboard of each student.