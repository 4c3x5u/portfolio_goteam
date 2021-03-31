from rest_framework.decorators import api_view
from rest_framework.response import Response
from ..serializers.board_serializer import CreateBoardSerializer


@api_view(['POST'])
def create_board(request):
    serializer = CreateBoardSerializer(data=request.data)
    if serializer.is_valid():
        return Response({
            'msg': 'Board created successfuly',
            'board_id': serializer.validated_data['board_id'],
            'team_id': serializer.validated_data['team_id']
        }, 201)
    return Response(serializer.errors, 400)
