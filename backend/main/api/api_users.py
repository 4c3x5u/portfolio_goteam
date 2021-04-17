from rest_framework.decorators import api_view
from rest_framework.response import Response
from ..models import User, Board
from ..util import validate_team_id, validate_board_id, authenticate


@api_view(['GET'])
def users(request):
    username = request.META.get('HTTP_AUTH_USER')
    token = request.META.get('HTTP_AUTH_TOKEN')

    authentication_response = authenticate(username, token)
    if authentication_response:
        return authentication_response

    team_id = request.query_params.get('team_id')
    validation_response = validate_team_id(team_id)
    if validation_response:
        return validation_response

    board_id = request.query_params.get('board_id')
    board, validation_response = validate_board_id(board_id)
    if validation_response:
        return validation_response

    users_list = User.objects.filter(team_id=team_id)
    board_users = User.objects.filter(board=board)

    return Response(list(map(
        lambda user: {
            'username': user.username,
            'isActive': user in board_users
        },
        users_list
    )), 200)
