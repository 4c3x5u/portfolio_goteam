from rest_framework.decorators import api_view
from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
import status

from ..serializers.boardserializer import BoardSerializer
from ..serializers.createboardserializer import CreateBoardSerializer
from ..serializers.deleteboardserializer import DeleteBoardSerializer
from ..models import Board
from ..validation.val_auth import authenticate, authorize, \
    not_authenticated_response
from ..validation.val_board import validate_board_id


class Boards(APIView):
    @staticmethod
    def post(request):
        serializer = CreateBoardSerializer(data={
            'team': request.data.get('team_id'),
            'name': request.data.get('name'),
            'auth_user': request.META.get('HTTP_AUTH_USER'),
            'auth_token': request.META.get('HTTP_AUTH_TOKEN')
        })
        if serializer.is_valid():
            serializer.save()
            return Response(serializer.data, status.HTTP_201_CREATED)
        return Response(serializer.errors, status.HTTP_400_BAD_REQUEST)

    @staticmethod
    def delete(request):
        board_id = request.query_params.get('id')
        serializer = DeleteBoardSerializer(data={
            'id': board_id or None,
            'auth_user': request.META.get('HTTP_AUTH_USER'),
            'auth_token': request.META.get('HTTP_AUTH_TOKEN')
        })
        if serializer.is_valid():
            serializer.delete()
            return Response(serializer.data, status.HTTP_200_OK)
        return Response(serializer.errors, status.HTTP_400_BAD_REQUEST)


@api_view(['POST', 'DELETE', 'PATCH'])
def boards(request):
    if request.method == 'DELETE':
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


