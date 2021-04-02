from rest_framework.test import APITestCase
from ..models import Task, Column, Board, Team


class UpdateTaskTests(APITestCase):
    def setUp(self):
        self.url = '/tasks/'
        self.task = Task.objects.create(
            title="Some Subtask Title",
            order=0,
            column=Column.objects.create(
                order=0,
                board=Board.objects.create(
                    team=Team.objects.create()
                )
            )
        )