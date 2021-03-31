from rest_framework.test import APITestCase
from ..models import Board, Team, User


class CreateBoardTests(APITestCase):
    url = '/board/'
    username = None
    team_id = None

    def setUp(self):
        team = Team.objects.create()
        self.team_id = team.id
        user = User.objects.create({
            'username': 'foooo',
            'password': 'barbarbar',
            'is_admin': True,
            'team': team
        })
        self.username = user.username
