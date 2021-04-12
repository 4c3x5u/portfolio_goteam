from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..serializers.ser_board import BoardSerializer
from ..serializers.ser_column import ColumnSerializer
from ..models import User, Board
from .util import authenticate, validate_team_id


@api_view(['POST', 'GET'])
def boards(request):
    auth_response = authenticate(request)
    if auth_response:
        return auth_response

    user = User.objects.get(username=request.META.get('HTTP_AUTH_USER'))

    if request.method == 'POST':
        # validate is_admin
        if not user.is_admin:
            return Response({
                'username': ErrorDetail(
                    string='Only the team admin can create a board.',
                    code='not_authorized'
                )
            }, 400)

        # validate team_id
        team_id = request.data.get('team_id')
        response = validate_team_id(team_id)
        if response:
            return response

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
        # validate team_id
        team_id = request.query_params.get('team_id')
        response = validate_team_id(team_id)
        if response:
            return response

        # create a board if none exists for the team
        team_boards = Board.objects.filter(team=team_id)
        if not team_boards:
            if user.is_admin:
                # create a board
                serializer = BoardSerializer(data={'team': team_id})
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
