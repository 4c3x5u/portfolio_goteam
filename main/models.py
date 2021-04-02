from django.db.models import *
import uuid


class Team(Model):
    invite_code = UUIDField(default=uuid.uuid4)


class User(Model):
    username = CharField(primary_key=True, max_length=35)
    password = BinaryField()
    is_admin = BooleanField(default=False)
    team = ForeignKey(Team, on_delete=CASCADE)


class Board(Model):
    team = ForeignKey(Team, on_delete=CASCADE)
    user = ManyToManyField(User)


class Column(Model):
    order = IntegerField()
    board = ForeignKey(Board, on_delete=CASCADE)


class Task(Model):
    title = CharField(max_length=50)
    description = TextField(null=True)
    # TODO: Make sure to delete this field if you find that you don't need it
    #       for React drag and drop controls
    order = IntegerField()
    column = ForeignKey(Column, on_delete=CASCADE)


# TODO: Add a 'done' bool field
class Subtask(Model):
    title = CharField(max_length=50)
    order = IntegerField()
    task = ForeignKey(Task, on_delete=CASCADE)
