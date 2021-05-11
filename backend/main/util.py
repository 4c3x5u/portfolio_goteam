import bcrypt
import json
import os
from rest_framework.serializers import ValidationError

from .models import User, Task, Subtask
from .serializers.board.ser_board import BoardSerializer
from .serializers.column.ser_column import ColumnSerializer


def create_admin(team, username_suffix=''):
    user = User.objects.create(
        username=f'teamadmin{username_suffix}',
        password=b'$2b$12$DKVJHUAQNZqIvoi.OMN6v.x1ZhscKhbzSxpOBMykHgTIMeeJpC6m'
                 b'e',
        is_admin=True,
        team=team
    )
    token = '$2b$12$yGUdlz0eMW3P6TAX07.CPuCA5u.t10uTEKCE2SQ5Vdm3VbnrHbpoK'
    return {'username': user.username,
            'password': user.password,
            'password_raw': 'barbarbar',
            'is_admin': user.is_admin,
            'team': user.team,
            'token': token if not username_suffix else bcrypt.hashpw(
                bytes(user.username, 'utf-8') + user.password,
                bcrypt.gensalt()
            ).decode('utf-8')}


def create_member(team, username_suffix=''):
    user = User.objects.create(
        username=f'teammember{username_suffix}',
        password=b'$2b$12$DKVJHUAQNZqIvoi.OMN6v.x1ZhscKhbzSxpOBMykHgTIMeeJpC6m'
                 b'e',
        is_admin=False,
        team=team
    )
    token = '$2b$12$qNhh2i1ZPU1qaIKncI7J6O2kr4XmuCWSwLEMJF653vyvDMIRekzLO'
    return {'username': user.username,
            'password': user.password,
            'password_raw': 'barbarbar',
            'is_admin': user.is_admin,
            'team': user.team,
            'token': token if not username_suffix else bcrypt.hashpw(
                bytes(user.username, 'utf-8') + user.password,
                bcrypt.gensalt()
            ).decode('utf-8')}


def create_board(name, team_id, team_admin):
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


def create_tutorial_tasks(user, column):
    path = os.path.abspath('main/data/tutorial_tasks.json')
    with open(path, 'r') as read_file:
        tutorial_tasks = json.load(read_file)

    # Subtasks cannot be created before the tasks are created, so two
    # iterations are needed. Otherwise, it would mean too many DB calls.
    tasks = [
        Task(
            title=task['title'],
            description=task['description'],
            order=i,
            column=column,
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
