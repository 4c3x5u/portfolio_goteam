from rest_framework.views import APIView
from rest_framework.response import Response
import status

from server.main.serializers.clientstate.ser_client_state import ClientStateSerializer


class ClientStateApiView(APIView):
    @staticmethod
    def get(request):
        serializer = ClientStateSerializer(data={
            'user': request.META.get('HTTP_AUTH_USER'),
            'auth_token': request.META.get('HTTP_AUTH_TOKEN'),
            'board_id': request.query_params.get('boardId') or 0
        })
        if not serializer.is_valid():
            return Response(serializer.errors, status.HTTP_400_BAD_REQUEST)
        return Response(serializer.data, status.HTTP_200_OK)
