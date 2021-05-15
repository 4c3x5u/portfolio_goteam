import json
import bcrypt

from main.models import User, Task, Subtask


class TutorialHelper:
    """
    Creates tutorial tasks inside the given column for the given user.
    """
    def __init__(self, user, column):
        self.user = user
        self.column = column

    def start(self):
        # CREATE A TEAM MEMBER
        User.objects.create(username=f'demo-member-{self.user.team_id}',
                            password=bcrypt.hashpw(b'securepassword',
                                                   bcrypt.gensalt()),
                            team_id=self.user.team_id)

        # CREATE TASKS
        with open('main/data/tutorial_tasks.json', 'r') as read_file:
            tutorial_tasks = json.load(read_file)

        # Subtasks cannot be created before the tasks are created, so two
        # iterations are needed. Otherwise, it would mean too many DB calls.
        tasks = [Task(title=task['title'],
                      description=task['description'],
                      order=i,
                      column=self.column,
                      user=self.user)
                 for i, task in enumerate(tutorial_tasks)]
        Task.objects.bulk_create(tasks)

        subtasks = [Subtask(title=title, task=tasks[ti], order=si)
                    for ti, task in enumerate(tutorial_tasks)
                    for si, title in enumerate(task['subtasks'])]
        Subtask.objects.bulk_create(subtasks)
