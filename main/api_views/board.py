from rest_framework.decorators import api_view
from rest_framework.response import Response
from main.serializers.board_serializer import BoardSerializer


@api_view(['POST'])
def board(request):
    serializer = BoardSerializer(data=request.data)
    if serializer.is_valid():
        new_board = serializer.save()
        return Response({
            'msg': 'Board created successfuly',
            'board_id': new_board.id,
            'team_id': new_board.team.id
        }, 201)
    return Response(serializer.errors, 400)
