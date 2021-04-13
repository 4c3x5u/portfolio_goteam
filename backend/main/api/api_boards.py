from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..serializers.ser_board import BoardSerializer
from ..serializers.ser_column import ColumnSerializer
from ..models import Board
from ..util import authenticate, authorize, validate_team_id


@api_view(['POST', 'GET'])
def boards(request):
    username = request.META.get('HTTP_AUTH_USER')
    token = request.META.get('HTTP_AUTH_TOKEN')

    authentication_response = authenticate(username, token)
    if authentication_response:
        return authentication_response

    if request.method == 'POST':
        authorization_response = authorize(username)
        if authorization_response:
            return authorization_response

        # validate team_id
        team_id = request.data.get('team_id')
        response = validate_team_id(team_id)
        if response:
            return response

        # create board
        board_serializer = BoardSerializer(
            data={'team': team_id,
                  'name': request.data.get('name')}
        )
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
        # validate team_id
        team_id = request.query_params.get('team_id')
        response = validate_team_id(team_id)
        if response:
            return response

        # create a board if none exists for the team
        team_boards = Board.objects.filter(team=team_id)
        if not team_boards:
            if not authorize(username):
                # create a board
                serializer = BoardSerializer(data={'team': team_id,
                                                   'name': 'New Board'})
                if not serializer.is_valid():
                    return Response({
                        'team_id': ErrorDetail(string='Invalid team ID.',
                                               code='invalid')
                    }, 400)
                board = serializer.save()

                # return a list containing only the new board
                return Response({
                    'boards': [
                        {'board_id': board.id, 'team_id': board.team.id}
                    ]
                }, 201)

            return Response({
                'team_id': ErrorDetail(string='Boards not found.',
                                       code='not_found')
            }, 404)

        # return pre-existing boards
        serializer = BoardSerializer(team_boards, many=True)
        return Response({'boards': serializer.data}, 200)
