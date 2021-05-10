from main.models import Team
from rest_framework.exceptions import ErrorDetail
from rest_framework.response import Response


# return (team, response)
def validate_team_id(team_id):
    if not team_id:
        return Response({
            'team_id': ErrorDetail(string='Team ID cannot be empty.',
                                   code='blank')
        }, 400)

    try:
        int(team_id)
    except ValueError:
        return Response({
            'team_id': ErrorDetail(string='Team ID must be a number.',
                                   code='invalid')
        }, 400)
