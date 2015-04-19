select
    s.stop_id,
    s.stop_name,
    s.stop_desc,
    s.stop_lat,
    s.stop_lon,
    s.location_type,
    111195 * st_distance(st_geomfromtext('point(%s %s)'), s.stop_geo) as stop_distance
from
    %s.stops s
where
    111195 * st_distance(st_geomfromtext('point(%s %s)'), s.stop_geo) < %s
order by
    stop_distance asc