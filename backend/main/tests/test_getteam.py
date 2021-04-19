from rest_framework.test import APITestCase
from ..models import Team
from ..util import new_admin


class GetTeamTests(APITestCase):
    endpoint = '/teams/?id='

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
