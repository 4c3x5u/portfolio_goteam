from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..serializers.ser_board import BoardSerializer
from ..serializers.ser_column import ColumnSerializer
from ..models import Board
from ..util import authenticate, authorize, validate_team_id


def create_board(team_id, name):
    # create board
    response = None
    board = None

    board_serializer = BoardSerializer(data={'team': team_id, 'name': name})
    if not board_serializer.is_valid():
        response = Response(board_serializer.errors, 400)
    else:
        board = board_serializer.save()

        # create four columns for the board
        for order in range(0, 4):
            column_serializer = ColumnSerializer(
                data={'board': board.id, 'order': order}
            )
            if not column_serializer.is_valid():
                response = Response(
                    column_serializer.errors, 400
                ) if not response else response
            column_serializer.save()

    return board, response


@api_view(['GET', 'POST', 'DELETE'])
def boards(request):
    username = request.META.get('HTTP_AUTH_USER')
    token = request.META.get('HTTP_AUTH_TOKEN')

    authentication_response = authenticate(username, token)
    if authentication_response:
        return authentication_response

    if request.method == 'GET':
        # validate team_id
        team_id = request.query_params.get('team_id')
        response = validate_team_id(team_id)
        if response:
            return response

        # create a board if none exists for the team
        team_boards = Board.objects.filter(team=team_id)
        if not team_boards:
            if not authorize(username):
                create_board_response = create_board(team_id, 'New Board')

                # return a list containing only the new board
                return Response({
                    'boards': [{'id': create_board_response.id, 'name': create_board_response.name}]
                }, 201)

            return Response({
                'team_id': ErrorDetail(string='Boards not found.',
                                       code='not_found')
            }, 404)

        # return pre-existing boards
        return Response({
            'boards': list(
                map(
                    lambda b: {'id': b['id'], 'name': b['name']},
                    BoardSerializer(team_boards, many=True).data
                )
            )
        }, 200)

    if request.method == 'POST':
        authorization_response = authorize(username)
        if authorization_response:
            return authorization_response

        # validate team_id
        team_id = request.data.get('team_id')
        response = validate_team_id(team_id)
        if response:
            return response

        board_name = request.data.get('name')

        board, create_board_response = create_board(team_id, board_name)
        if create_board_response:
            return create_board_response

        # return success response
        return Response({
            'msg': 'Board creation successful.',
            'id': board.id,
        }, 201)

    if request.method == 'DELETE':
        authorization_response = authorize(username)
        if authorization_response:
            return authorization_response

        board_id = request.query_params.get('id')

        if not board_id:
            return Response({
                'board_id': ErrorDetail(string='Board ID cannot be empty.',
                                        code='blank')
            }, 400)

        try:
            int(board_id)
        except ValueError:
            return Response({
                'board_id': ErrorDetail(string='Board ID must be a number.',
                                        code='invalid')
            }, 400)

        try:
            Board.objects.get(id=board_id).delete()
        except Board.DoesNotExist:
            return Response({
                'board_id': ErrorDetail(string='Board not found.',
                                        code='not_found')
            }, 404)

        return Response({
            'msg': 'Board deleted successfully.',
            'id': board_id,
        })
