from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..serializers.ser_board import BoardSerializer
from ..serializers.ser_column import ColumnSerializer
from ..models import Board, Team, User
import bcrypt


def validate_username(username):
    if not username:
        return Response({
            'username': ErrorDetail(string="Username cannot be empty.",
                                    code='blank')
        }, 400)
    try:
        User.objects.get(username=username)
    except User.DoesNotExist:
        return Response({
            'username': ErrorDetail(string="Invalid username.",
                                    code='invalid')
        }, 400)


def validate_team_id(id):
    if not id:
        return Response({
            'team_id': ErrorDetail(string='Team ID cannot be empty.',
                                   code='null')
        }, 400)
    try:
        Team.objects.get(id=id)
    except Team.DoesNotExist:
        return Response({
            'team_id': ErrorDetail(string='Team not found.',
                                   code='not_found')
        }, 404)

def validate_authentication_token(token, username, password):
    if not token:
        return Response({
            'token': ErrorDetail(
                string='Authentication token cannot be empty.',
                code='blank'
            )
        }, 400)
    no_match_response = Response({
        'token': ErrorDetail(string='Invalid authentication token.',
                             code='invalid')
    }, 400)
    try:
        tokens_match = bcrypt.checkpw(
            bytes(username, 'utf-8') + password,
            bytes(token, 'utf-8'))
        if not tokens_match:
            return no_match_response
    except ValueError:
        return no_match_response


@api_view(['POST', 'GET'])
def boards(request):
    if request.method == 'POST':
        username = request.META.get('HTTP_AUTH_USER')
        validate_username(username)
        user = User.objects.get(username=username)

        # validate is_admin
        if not user.is_admin:
            return Response({
                'username': ErrorDetail(
                    string='Only the team admin can create a board.',
                    code='not_authorized'
                )
            }, 400)

        team_id = request.data.get('team_id')
        validate_team_id(team_id)

        token = request.META.get('HTTP_AUTH_TOKEN')
        validate_authentication_token(token, user.username, user.password)

        # create board
        board_serializer = BoardSerializer(data={'team': team_id})
        if not board_serializer.is_valid():
            return Response(board_serializer.errors, 400)
        board = board_serializer.save()

        # create four columns for the board
        for order in range(0, 4):
            column_serializer = ColumnSerializer(
                data={'board': board.id, 'order': order}
            )
            if not column_serializer.is_valid():
                return Response(column_serializer.errors, 400)
            column_serializer.save()

        # return success response
        return Response({
            'msg': 'Board creation successful.',
            'board_id': board.id
        }, 201)

    if request.method == 'GET':
        # validate username
        username = request.data.get('username')
        validate_username(username)
        user = User.objects.get(username)

        # validate team_id
        team_id = request.query_params.get('team_id')
        validate_team_id(team_id)

        # TODO: validate authentication token
        token = request.query_params.get('token')
        validate_authentication_token(token, user.username, user.password)

        # create a board if none exists for the team
        team_boards = Board.objects.filter(team=team_id)
        if not team_boards:
            serializer = BoardSerializer(data={'team': team_id})
            if not serializer.is_valid():
                return Response({
                    'team_id': ErrorDetail(string='Invalid team ID.',
                                           code='invalid')
                }, 404)
            board = serializer.save()

            # return the new board
            return Response({
                'boards': [{'board_id': board.id, 'team_id': board.team.id}]
            }, 201)

        # return pre-existing boards
        return Response({'boards': BoardSerializer(team_boards, many=True)},
                        200)
