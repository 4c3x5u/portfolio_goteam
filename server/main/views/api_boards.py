from rest_framework.views import APIView
from rest_framework.response import Response
import status

from server.main.serializers.board.ser_board_create import CreateBoardSerializer
from server.main.serializers.board.ser_board_update import UpdateBoardSerializer
from server.main.serializers.board.ser_board_delete import DeleteBoardSerializer


class BoardsApiView(APIView):
    @staticmethod
    def post(request):
        serializer = CreateBoardSerializer(
            data={'team': request.data.get('team_id'),
                  'name': request.data.get('name'),
                  'auth_user': request.META.get('HTTP_AUTH_USER'),
                  'auth_token': request.META.get('HTTP_AUTH_TOKEN')}
        )
        if serializer.is_valid():
            serializer.save()
            return Response(serializer.data, status.HTTP_201_CREATED)
        return Response(serializer.errors, status.HTTP_400_BAD_REQUEST)

    @staticmethod
    def patch(request):
        serializer = UpdateBoardSerializer(
            data={'board': request.query_params.get('id') or None,
                  'payload': request.data,
                  'auth_user': request.META.get('HTTP_AUTH_USER'),
                  'auth_token': request.META.get('HTTP_AUTH_TOKEN')}
        )
        if serializer.is_valid():
            serializer.save()
            return Response(serializer.data, status.HTTP_200_OK)
        return Response(serializer.errors, status.HTTP_400_BAD_REQUEST)

    @staticmethod
    def delete(request):
        serializer = DeleteBoardSerializer(
            data={'board': request.query_params.get('id') or None,
                  'auth_user': request.META.get('HTTP_AUTH_USER'),
                  'auth_token': request.META.get('HTTP_AUTH_TOKEN')}
        )
        if serializer.is_valid():
            serializer.delete()
            return Response(serializer.data, status.HTTP_200_OK)
        return Response(serializer.errors, status.HTTP_400_BAD_REQUEST)
