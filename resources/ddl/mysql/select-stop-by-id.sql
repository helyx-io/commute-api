select
    s.stop_id,
    s.stop_name,
    s.stop_desc,
    s.stop_lat,
    s.stop_lon,
    s.location_type
from
    %s.stops s
where
    s.stop_id = %s
