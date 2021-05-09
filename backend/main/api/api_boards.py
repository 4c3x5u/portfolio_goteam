from rest_framework.decorators import api_view
from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
import status

from ..serializers.boardserializer import BoardSerializer
from ..serializers.createboardserializer import CreateBoardSerializer
from ..models import Board, Team
from ..validation.val_auth import authenticate, authorize, \
    not_authenticated_response
from ..validation.val_team import validate_team_id
from ..validation.val_board import validate_board_id
from ..util import create_board


class Boards(APIView):
    @staticmethod
    def post(request):
        serializer = CreateBoardSerializer(data={
            'auth_user': request.META.get('HTTP_AUTH_USER'),
            'auth_token': request.META.get('HTTP_AUTH_TOKEN'),
            'team': request.data.get('team_id'),
            'name': request.data.get('name')
        })
        if not serializer.is_valid():
            return Response(serializer.errors, status.HTTP_400_BAD_REQUEST)
        serializer.save()
        return Response(serializer.data, status.HTTP_201_CREATED)


@api_view(['POST', 'DELETE', 'PATCH'])
def boards(request):
    username = request.META.get('HTTP_AUTH_USER')
    token = request.META.get('HTTP_AUTH_TOKEN')

    user, authentication_response = authenticate(username, token)
    if authentication_response:
        return authentication_response

    if request.method == 'POST':
        authorization_response = authorize(username)
        if authorization_response:
            return authorization_response

        # validate team_id
        team_id = request.data.get('team_id')
        validation_response = validate_team_id(team_id)
        if validation_response:
            return validation_response

        try:
            team = Team.objects.prefetch_related('user_set').get(id=team_id)
        except Team.DoesNotExist:
            return Response({
                'team_id': ErrorDetail(string='Team not found.',
                                       code='not_found')
            }, 404)

        if team.id != user.team_id:
            return not_authenticated_response

        board_name = request.data.get('name')

        board, error_response = create_board(
            name=board_name,
            team_id=team.id,
            team_admin=team.user_set.get(username=username)
        )
        if error_response:
            return error_response

        return Response({
            'msg': 'Board creation successful.',
            'id': board.id,
        }, 201)

    if request.method == 'DELETE':
        authorization_response = authorize(username)
        if authorization_response:
            return authorization_response

        board_id = request.query_params.get('id')

        validation_response = validate_board_id(board_id)
        if validation_response:
            return validation_response

        try:
            board = Board.objects.get(id=board_id)
        except Board.DoesNotExist:
            return Response({
                'board_id': ErrorDetail(string='Board not found.',
                                        code='not_found')
            }, 404)

        if board.team_id != user.team_id:
            return not_authenticated_response

        board.delete()

        return Response({
            'msg': 'Board deleted successfully.',
            'id': board_id,
        })

    if request.method == 'PATCH':
        authorization_response = authorize(username)
        if authorization_response:
            return authorization_response

        board_id = request.query_params.get('id')
        validation_response = validate_board_id(board_id)
        if validation_response:
            return validation_response

        try:
            board = Board.objects.get(id=board_id)
        except Board.DoesNotExist:
            return Response({
                'board_id': ErrorDetail(string='Board not found.',
                                        code='not_found')
            }, 404)

        if board.team_id != user.team_id:
            return not_authenticated_response

        serializer = BoardSerializer(board, data=request.data, partial=True)
        if not serializer.is_valid():
            return Response(serializer.errors, 400)
        serializer.save()

        return Response({
            'msg': 'Board updated successfuly.',
            'id': serializer.data['id'],
        }, 200)


