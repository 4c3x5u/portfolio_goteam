from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..models import Team
from ..validation.val_auth import \
    authenticate, authorize, not_authenticated_response
from ..validation.val_team import validate_team_id


# not in use – maintained for demonstration purposes
@api_view(['GET'])
def teams(request):
    auth_user = request.META.get('HTTP_AUTH_USER')
    auth_token = request.META.get('HTTP_AUTH_TOKEN')

    user, authentication_response = authenticate(auth_user, auth_token)
    if authentication_response:
        return authentication_response

    authorization_response = authorize(auth_user)
    if authorization_response:
        return authorization_response

    team_id = request.query_params.get('team_id')

    validation_response = validate_team_id(team_id)
    if validation_response:
        return validation_response

    try:
        team = Team.objects.get(id=team_id)
    except Team.DoesNotExist:
        return Response({
            'team_id': ErrorDetail(string='Team not found.',
                                   code='not_found')
        }, 404)

    if team.id != user.team_id:
        return not_authenticated_response

    return Response({
        'id': team.id,
        'inviteCode': team.invite_code
    }, 200)


