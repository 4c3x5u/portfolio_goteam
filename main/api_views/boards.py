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
                'msg': 'Board created successfuly',
                'board_id': board_response.id,
                'team_id': board_response.team.id
            }, 201)
        return Response(serializer.errors, 400)

    if request.method == 'GET':
        team_id = request.query_params.get('team_id')
        boards = Board.objects.filter(team=team_id)
        if not boards:
            return Response({
                'team_id': ErrorDetail(
                    string='No boards found for this team ID.',
                    code='not_found'
                )
            })
        serializer = BoardSerializer(
            data=list(map(lambda b: b.__dict__, boards)),
            many=True
        )
        if serializer.is_valid():
            return Response({
                'boards': serializer.validated_data
            }, 200)
        return Response(serializer.errors, 400)
