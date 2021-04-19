from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Team
from ..util import new_admin


class GetTeamTests(APITestCase):
    endpoint = '/teams/?team_id='

    def setUp(self):
        self.team = Team.objects.create()
        self.member = new_admin(self.team)

    def test_success(self):
        response = self.client.get(f'{self.endpoint}{self.team.id}')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {
            'id': self.team.id,
            'inviteCode': self.team.invite_code
        })

    def test_team_id_empty(self):
        response = self.client.get(self.endpoint)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team ID cannot be empty.',
                                   code='blank')
        })

    def test_team_id_invalid(self):
        response = self.client.get(f'{self.endpoint}asdfsa')
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team ID must be a number.',
                                   code='invalid')
        })


    def test_team_not_found(self):
        response = self.client.get(f'{self.endpoint}12314241')
        self.assertEqual(response.status_code, 404)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team not found.',
                                   code='not_found')
        })

