from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..models import User
from ..util import validate_team_id, authenticate


@api_view(['GET'])
def users(request):
    username = request.META.get('HTTP_AUTH_USER')
    token = request.META.get('HTTP_AUTH_TOKEN')

    authentication_response = authenticate(username, token)
    if authentication_response:
        return authentication_response

    team_id = request.query_params.get('team_id')

    validation_response = validate_team_id(team_id)
    if validation_response:
        return validation_response

    users_queryset = User.objects.filter(team_id=team_id)
    return Response(list(map(
        lambda user: {'username': user.username},
        users_queryset
    )), 200)
