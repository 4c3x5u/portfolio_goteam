from rest_framework.serializers import ValidationError
from main.models import User, Task, Subtask
from main.serializers.board.ser_board import BoardSerializer
from main.serializers.column.ser_column import ColumnSerializer
import bcrypt


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


def create_board(name, team_id, team_admin):  # -> (board, response)
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
    tasks = [
        Task(title='Drag and Drop Controls',
             description='Complete this task to gain familiarity with the '
                         'drag and drop controls.',
             order=0,
             column=column,
             user=user),
        Task(title='Creating Tasks',
             description='Complete this task to gain familiarity with the '
                         'task creation process.',
             order=1,
             column=column,
             user=user),
        Task(title='Editing Tasks',
             description='Complete this task to gain familiarity with the '
                         'edit task process.',
             order=2,
             column=column,
             user=user),
    ]
    Task.objects.bulk_create(tasks)

    subtasks = [
        # 'Drag and Drop Controls' subtasks
        Subtask(title='Drag and drop this task to the GO column.',
                order=0,
                task=tasks[0]),
        Subtask(title='Drag and drop this task to the DONE column.',
                order=1,
                task=tasks[0]),

        # 'Creating Tasks' subtasks
        Subtask(title='Move this task to the GO column.',
                order=0,
                task=tasks[1]),
        Subtask(title='Activate the task creation dialogue by clicking the '
                      'plus button inside the INBOX column.',
                order=1,
                task=tasks[1]),
        Subtask(title='Give the task a title.', order=2, task=tasks[1]),
        Subtask(title='Give the task a description.', order=3, task=tasks[1]),
        Subtask(title='Add a couple of subtasks to the task.',
                order=4,
                task=tasks[1]),
        Subtask(title='Click the CREATE button.', order=5, task=tasks[1]),
        Subtask(title='Move this task to the DONE column.',
                order=6,
                task=tasks[1]),

        # 'Editing Tasks' subtasks
        Subtask(title='Move this task to the GO column.',
                order=0,
                task=tasks[2]),
        Subtask(title='Activate the edit task dialogue by right-clicking the '
                      'task you just created, and then clicking EDIT.',
                order=1,
                task=tasks[2]),
        Subtask(title='Edit the details of the task — edit its title and '
                      'description. Add and remove subtasks from it.',
                order=2,
                task=tasks[2]),
        Subtask(title='Click the SUBMIT button.',
                order=3,
                task=tasks[2]),
        Subtask(title='You know what to do with this task. :)',
                order=4,
                task=tasks[2]),
    ]
    Subtask.objects.bulk_create(subtasks)


