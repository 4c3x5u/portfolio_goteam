from main.models import Team, User
from rest_framework.exceptions import ErrorDetail
from rest_framework.response import Response
from .serializers.ser_board import BoardSerializer
from .serializers.ser_column import ColumnSerializer
from .models import Board, Column, Task
import bcrypt


def new_admin(team, username_suffix=''):
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


def new_member(team, username_suffix=''):
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
            ).decode('utf-8')
}


not_authenticated_response = Response({
    'auth': ErrorDetail(string="Authentication failure.",
                        code='not_authenticated')
}, 403)


def authenticate(username, token):  # -> (team id, response)
    try:
        user = User.objects.get(username=username)
    except (User.DoesNotExist, ValueError):
        return None, not_authenticated_response

    try:
        tokens_match = bcrypt.checkpw(
            bytes(user.username, 'utf-8') + user.password,
            bytes(token, 'utf-8'))
        if not tokens_match:
            return None, not_authenticated_response
    except ValueError:
        return None, not_authenticated_response

    return user.team_id, None


not_authorized_response = Response({
    'auth': ErrorDetail(string='You must be an admin to do this.',
                        code='not_authorized')
}, 403)


def authorize(username):
    try:
        user = User.objects.get(username=username)
        if not user.is_admin:
            return not_authorized_response
    except User.DoesNotExist:
        return not_authorized_response


def validate_team_id(team_id):  # -> (team, validation_response)
    if not team_id:
        return None, Response({
            'team_id': ErrorDetail(string='Team ID cannot be empty.',
                                   code='blank')
        }, 400)

    try:
        int(team_id)
    except ValueError:
        return None, Response({
            'team_id': ErrorDetail(string='Team ID must be a number.',
                                   code='invalid')
        }, 400)

    try:
        team = Team.objects.get(id=team_id)
    except Team.DoesNotExist:
        return None, Response({
            'team_id': ErrorDetail(string='Team not found.',
                                   code='not_found')
        }, 404)

    return team, None


def create_board(team_id, name):  # -> (board, response)
    board_serializer = BoardSerializer(data={'team': team_id, 'name': name})
    if not board_serializer.is_valid():
        return None, Response(board_serializer.errors, 400)

    board = board_serializer.save()

    # create four columns for the board
    for order in range(0, 4):
        column_serializer = ColumnSerializer(
            data={'board': board.id, 'order': order}
        )
        if not column_serializer.is_valid():
            return board, Response(
                column_serializer.errors, 400
            )
        column_serializer.save()

    return board, None


def validate_board_id(board_id):  # -> (board, response)
    if not board_id:
        return None, Response({
            'board_id': ErrorDetail(string='Board ID cannot be empty.',
                                    code='blank')
        }, 400)

    try:
        int(board_id)
    except ValueError:
        return None, Response({
            'board_id': ErrorDetail(string='Board ID must be a number.',
                                    code='invalid')
        }, 400)

    try:
        board = Board.objects.get(id=board_id)
    except Board.DoesNotExist:
        return None, Response({
            'board_id': ErrorDetail(string='Board not found.',
                                    code='not_found')
        }, 404)

    return board, None


def validate_column_id(column_id):  # -> (column, response)
    if not column_id:
        return None, Response({
            'column_id': ErrorDetail(string='Column ID cannot be empty.',
                                     code='blank')
        }, 400)

    try:
        int(column_id)
    except ValueError:
        return None, Response({
            'column_id': ErrorDetail(string='Column ID must be a number.',
                                     code='invalid')
        }, 400)

    try:
        column = Column.objects.get(id=column_id)
    except Column.DoesNotExist:
        return None, Response({
            'column_id': ErrorDetail(string='Column not found.',
                                     code='not_found')
        }, 404)

    return column, None


def validate_task_id(task_id):  # -> (task, response)
    if not task_id:
        return None, Response({
            'task_id': ErrorDetail(string='Task ID cannot be empty.',
                                   code='blank')
        }, 400)

    try:
        int(task_id)
    except ValueError:
        return None, Response({
            'task_id': ErrorDetail(string='Task ID must be a number.',
                                   code='invalid')
        }, 400)

    try:
        task = Task.objects.get(id=task_id)
    except Task.DoesNotExist:
        return None, Response({
            'task_id': ErrorDetail(string='Task not found.',
                                   code='not_found')
        }, 404)

    return task, None


def validate_is_active(is_active):
    is_empty_response = Response({
        'is_active': ErrorDetail(string='Is Active cannot be empty.',
                                 code='blank')
    }, 400)

    try:
        if not str(is_active):
            return None, is_empty_response
    except ValueError:
        return None, is_empty_response

    if not isinstance(is_active, bool):
        return None, Response({
            'is_active': ErrorDetail(string='Is Valid must be a boolean.',
                                     code='invalid')
        }, 400)

    return is_active, None
