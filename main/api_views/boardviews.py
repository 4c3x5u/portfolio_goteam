from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..serializers.boardserializer import BoardSerializer
from ..models import Board, Team, User


@api_view(['POST', 'GET'])
def boards(request):
    if request.method == 'POST':
        username = request.data.get('username')
        if not username:
            error = ErrorDetail(string="Username cannot be empty.",
                                code='blank')
            return Response({'username': error}, 400)
        try:
            user = User.objects.get(username=username)
        except User.DoesNotExist:
            error = ErrorDetail(string="Invalid username.", code='invalid')
            return Response({'username': error}, 400)
        if not user.is_admin:
            error = ErrorDetail(
                string='Only the team admin can create a board.',
                code='not_authorized'
            )
            return Response({'username': error}, 400)
        team_id = request.data.get('team_id')
        if not team_id:
            error = ErrorDetail(string='Team ID cannot be empty.',
                                code='blank')
            return Response({'team_id': error}, 400)
        try:
            Team.objects.get(id=team_id)
        except Team.DoesNotExist:
            error = ErrorDetail(string='Team not found.', code='not_found')
            return Response({'team_id': error}, 404)
        serializer = BoardSerializer(data={'team': team_id})
        if not serializer.is_valid():
            return Response(serializer.errors, 400)
        board = serializer.save()
        return Response(
            {'board_id': board.id, 'team_id': board.team.id},
            201
        )

    if request.method == 'GET':
        team_id = request.query_params.get('team_id')
        if not team_id:
            error = ErrorDetail(string='Team ID cannot be empty.', code='null')
            return Response({'team_id': error}, 400)
        try:
            Team.objects.get(id=team_id)
        except Team.DoesNotExist:
            error = ErrorDetail(string='Team not found.', code='not_found')
            return Response({'team_id': error}, 404)
        team_boards = Board.objects.filter(team=team_id)
        if not team_boards:
            # TODO: Create a board if not found and don't return error
            error = ErrorDetail(string='No boards found for this team.',
                                code='not_found')
            return Response({'team_id': error}, 404)
        serializer = BoardSerializer(team_boards, many=True)
        return Response({'boards': serializer.data}, 200)
