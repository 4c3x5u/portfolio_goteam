from django.db.models import *
import uuid


class Team(Model):
    id = AutoField(primary_key=True, db_index=True)
    invite_code = UUIDField(default=uuid.uuid4)


class User(Model):
    username = CharField(primary_key=True, max_length=35, db_index=True)
    password = BinaryField()
    is_admin = BooleanField(default=False)
    team = ForeignKey(Team, on_delete=CASCADE)


class Board(Model):
    id = AutoField(primary_key=True, db_index=True)
    name = CharField(max_length=35)
    team = ForeignKey(Team, on_delete=CASCADE, db_index=True)
    user = ManyToManyField(User, db_index=True)


class Column(Model):
    id = AutoField(primary_key=True, db_index=True)
    order = IntegerField()
    board = ForeignKey(Board, on_delete=CASCADE, db_index=True)


class Task(Model):
    id = AutoField(primary_key=True, db_index=True)
    title = CharField(max_length=50)
    description = TextField(blank=True, null=True)
    order = IntegerField()
    column = ForeignKey(Column, on_delete=CASCADE, db_index=True)
    user = ForeignKey(User, null=True, on_delete=SET_NULL)


class Subtask(Model):
    id = AutoField(primary_key=True, db_index=True)
    title = CharField(max_length=50)
    order = IntegerField()
    task = ForeignKey(Task, on_delete=CASCADE, db_index=True)
    done = BooleanField(default=False)
