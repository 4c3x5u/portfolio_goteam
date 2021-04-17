from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..models import User


@api_view(['GET'])
def users(request):
    team_id = request.query_params.get('team_id')
    if not team_id:
        return Response({
            'team_id': ErrorDetail(string='Team ID cannot be empty.',
                                   code='blank')
        }, 400)

    users_queryset = User.objects.filter(team_id=team_id)
    return Response(list(map(
        lambda user: {'username': user.username},
        users_queryset
    )), 200)
