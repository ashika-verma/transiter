from sqlalchemy import Column, Integer, String, ForeignKey, Boolean, Numeric, UniqueConstraint
from sqlalchemy.orm import relationship, backref

from .base import Base


class Stop(Base):
    __tablename__ = 'stop'

    pk = Column(Integer, primary_key=True)
    id = Column(String)
    system_id = Column(String, ForeignKey('system.id'))
    parent_stop_pk = Column(Integer, ForeignKey('stop.pk'), index=True)
    #TODO: remove station_pk and station relationship, make stops instead
    station_pk = Column(Integer, ForeignKey('station.pk'), index=True)

    name = Column(String)
    longitude = Column(Numeric(precision=9, scale=6))
    latitude = Column(Numeric(precision=9, scale=6))
    is_station = Column(Boolean)

    system = relationship(
        'System',
        back_populates='stops')
    # NOTE: this relationship is a little tricky and this definition follows
    # https://docs.sqlalchemy.org/en/latest/_modules/examples/adjacency_list/adjacency_list.html
    child_stops = relationship(
        'Stop',
        cascade='all, delete-orphan',
        backref=backref('parent_stop', remote_side=pk),
    )
    direction_name_rules = relationship(
        'DirectionNameRule',
        back_populates='stop',
        cascade='all, delete-orphan',
        order_by='DirectionNameRule.priority')
    stop_events = relationship(
        'StopTimeUpdate',
        back_populates='stop',
        cascade='all, delete-orphan')
    service_pattern_vertices = relationship(
        'ServicePatternVertex',
        back_populates='stop',
        cascade='all, delete-orphan')
    stop_id_aliases = relationship(
        'StopIdAlias',
        back_populates='stop',
        cascade='all, delete-orphan')

    __table_args__ = (UniqueConstraint('system_id', 'id'), )

    _short_repr_list = ['id', 'system_id', 'name']
