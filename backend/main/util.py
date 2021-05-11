from rest_framework.serializers import ValidationError
from main.models import User
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

