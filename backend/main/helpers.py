from rest_framework.serializers import ValidationError
import bcrypt
import json
import shortuuid

from .models import User, Task, Subtask
from .serializers.board.ser_board import BoardSerializer
from .serializers.column.ser_column import ColumnSerializer


class UserHelper:
    def __init__(self, team):
        self.team = team
        self._unique_identifier = shortuuid.uuid()
        self._user_counter = 0

    def create(self, is_admin=False):
        """
        Creates a new user and returns a dictionary containing user data, as
        well as the token.
        """
        user = User.objects.create(
            username=f'{self._unique_identifier}-{self._user_counter}',
            password=b'$2b$12$DKVJHUAQNZqIvoi.OMN6v.x1ZhscKhbzSxpOBMykHgTIMeeJ'
                     b'pC6me',
            is_admin=is_admin,
            team=self.team
        )
        self._user_counter += 1
        return {'username': user.username,
                'password': user.password,
                'password_raw': 'barbarbar',
                'is_admin': user.is_admin,
                'team': user.team,
                'token': bcrypt.hashpw(
                    bytes(user.username, 'utf-8') + user.password,
                    bcrypt.gensalt()
                ).decode('utf-8')}


class BoardHelper:
    @staticmethod
    def create(name, team_id, team_admin):
        """
        Creates a board, and four columns for it.
        """
        board_serializer = BoardSerializer(data={'team': team_id, 'name': name})
        if not board_serializer.is_valid():
            raise ValidationError({'boards': board_serializer.errors})
        board = board_serializer.save()

        board.user.add(team_admin)

        # create four columns for the board
        column_data = [
            {'board': board.id, 'order': order} for order in range(0, 4)
        ]

        column_serializer = ColumnSerializer(data=column_data, many=True)
        if not column_serializer.is_valid():
            raise ValidationError({'columns': column_serializer.errors})

        column_serializer.save()

        return board


class TutorialHelper:
    @staticmethod
    def initiate(user, ready_column):
        """
        Creates tutorial objects for a newly registered admin to go through.
        """

        # CREATE A TEAM MEMBER
        User.objects.create(username=f'demo-member-{user.team_id}',
                            password=b'securepassword',
                            team_id=user.team_id)

        # CREATE TASKS
        with open('main/data/tutorial_tasks.json', 'r') as read_file:
            tutorial_tasks = json.load(read_file)

        # Subtasks cannot be created before the tasks are created, so two
        # iterations are needed. Otherwise, it would mean too many DB calls.
        tasks = [
            Task(
                title=task['title'],
                description=task['description'],
                order=i,
                column=ready_column,
                user=user
            )
            for i, task in enumerate(tutorial_tasks)
        ]
        Task.objects.bulk_create(tasks)

        subtasks = [
            Subtask(title=title, task=tasks[ti], order=si)
            for ti, task in enumerate(tutorial_tasks)
            for si, title in enumerate(task['subtasks'])
        ]
        Subtask.objects.bulk_create(subtasks)
