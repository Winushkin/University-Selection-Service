SELECT u.*,
       r.name,
       s.cost,
       s.budget_points,
       s.contract_points
FROM universities.universities u
JOIN universities.regions r ON u.region_id = r.id
JOIN universities.specialities s ON u.id = s.university_id
WHERE s.name = $1;