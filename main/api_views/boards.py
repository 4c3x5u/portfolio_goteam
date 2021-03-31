from rest_framework.decorators import api_view
from rest_framework.response import Response
from main.serializers.createboardserializer import CreateBoardSerializer
from main.serializers.listboardsserializer import ListBoardsSerializer
from main.models import Board


@api_view(['POST', 'GET'])
def board(request):
    if request.method == 'POST':
        serializer = CreateBoardSerializer(data=request.data)
        if serializer.is_valid():
            new_board = serializer.save()
            return Response({
                'msg': 'Board created successfuly',
                'board_id': new_board.id,
                'team_id': new_board.team.id
            }, 201)
        return Response(serializer.errors, 400)

    if request.method == 'GET':
        serializer = ListBoardsSerializer(
            data={'team_id': request.query_params.get('team_id')}
        )
        if serializer.is_valid():
            return Response({
                'boards': serializer.get_list(serializer.validated_data['team_id']),
            }, 200)
        return Response(serializer.errors, 400)
