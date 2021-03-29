from django.db.models import *


class Team(Model):
    pass


class User(Model):
    username = CharField(primary_key=True, max_length=35)
    password = CharField(max_length=255)
    is_admin = BooleanField(default=False)
    team = ForeignKey(Team, on_delete=CASCADE)


class Board(Model):
    id = AutoField(primary_key=True)
    team = ForeignKey(Team, on_delete=CASCADE)
    user = ManyToManyField(User)


class Column(Model):
    order = IntegerField()
    board = ForeignKey(Board, on_delete=CASCADE)


class Task(Model):
    title = CharField(max_length=50)
    description = TextField(null=True)
    order = IntegerField()
    column = ForeignKey(Column, on_delete=CASCADE)


class Subtask(Model):
    title = CharField(max_length=50)
    order = IntegerField()
    task = ForeignKey(Task, on_delete=CASCADE)
