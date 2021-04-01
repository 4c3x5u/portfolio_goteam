from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..serializers.boardserializer import BoardSerializer
from ..models import Board


@api_view(['POST', 'GET'])
def boards(request):
    if request.method == 'POST':
        serializer = BoardSerializer(data={
            'username': request.data.get('username'),
            'team_id': request.data.get('team_id')
        })
        if serializer.is_valid():
            board_response = serializer.save()
            return Response({
                'board_id': board_response.id,
                'team_id': board_response.team.id
            }, 201)
        return Response(serializer.errors, 400)

    if request.method == 'GET':
        team_id = request.query_params.get('team_id')
        print(f'§TEAM_ID:{team_id}')
        if not team_id:
            return Response({
                'team_id': ErrorDetail(string='Team ID cannot be empty.',
                                       code='null')
            }, 400)
        team_boards = Board.objects.filter(team=team_id)
        if not team_boards:
            return Response({
                'team_id': ErrorDetail(string='No boards found.',
                                       code='not_found')
            }, 404)
        return Response({
            'boards': BoardSerializer(team_boards, many=True).data
        }, 200)
