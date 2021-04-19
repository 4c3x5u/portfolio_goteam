from rest_framework.decorators import api_view
from rest_framework.response import Response
from ..util import validate_team_id


@api_view(['GET'])
def teams(request):
    team_id = request.query_params.get('team_id')
    team, validation_response = validate_team_id(team_id)
    if validation_response:
        return validation_response
    return Response({
        'id': team.id,
        'inviteCode': team.invite_code
    }, 200)


