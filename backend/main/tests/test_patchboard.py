from rest_framework.test import APITestCase
from ..models import Board, Team
from ..util import new_admin


class PatchBoardTests(APITestCase):
    endpoint = '/boards/?id='

    def setUp(self):
        team = Team.objects.create()
        self.admin = new_admin(team)
        self.board = Board.objects.create(name='Some Board',
                                          team=team)

    def test_success(self):
        response = self.client.patch(f'{self.endpoint}{self.board.id}',
                                     {'name': 'New Title'},
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {
            'msg': 'Board updated successfuly.',
            'id': self.board.id,
        })
        self.assertEqual(Board.objects.get(id=self.board.id).name,
                         'New Title')
